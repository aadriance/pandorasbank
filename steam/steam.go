package steam

import (
	"os"
	"path/filepath"
	"github.com/andygrunwald/vdf"
)

const (
	ELESTRALS_AWAKENED_PLAYTEST = "4168870"
)


func GetDefaultGamePath(customGamePath string) string {
	if isValidPath(customGamePath) {
		return customGamePath
	}

	// Attempt to locate os-based steam installation
	steamPath := getSteamInstallDirectory()
	if !isValidPath(steamPath) {
		return ""
	}

	// Acquire libraryfolders.vdf - defining which library on the filesystem each game lives in
	vdfPath := filepath.Join(steamPath, STEAM_LIBRARIES_VDF)
	if !isValidPath(vdfPath) {
		return ""
	}

	libraryPath := getDesiredLibraryPath(vdfPath)
	if !isValidPath(libraryPath) {
		return ""
	}

	// TODO: Does app have a different extension (no extension) on Mac/Linux?  Maybe move to os-specific const
	appPath := filepath.Join(libraryPath, "steamapps", "common", "Elestrals Awakened Playtest", "ElestralsAwakened-Playtest.exe")
	if !isValidPath(appPath) {
		return ""
	}

	return appPath
}

func getDesiredLibraryPath(vdfPath string) string {
	vdfFile, err := os.Open(vdfPath)
	if err != nil {
		return ""
	}

	parser := vdf.NewParser(vdfFile)
	content, err := parser.Parse()
	if err != nil {
		return ""
	}

	libraryMap := getChildMap(content, "libraryfolders")
	for _, library := range libraryMap {
		if library != nil {
			if library, ok := library.(map[string]interface{}); ok {
				if appRecord := getChildValue(library, "apps", ELESTRALS_AWAKENED_PLAYTEST); appRecord != "" {
				  // Found the Elestrals Awakened Playtest app in this library - use as base path
					return getChildValue(library, "path")
				}
			}
		}
	}
	return ""
}

func getChildValue(input map[string]interface{}, keys ...string) string {
	return getContent[string](input, "", keys...)
}

func getChildMap(input map[string]interface{}, keys ...string) map[string]interface{} {
	return getContent[map[string]interface{}](input, nil, keys...)
}


func getContent[T map[string]interface{} | string](input map[string]interface{}, otherwise T, keys ...string) T {
	if len(keys) == 0 {
		return otherwise 
	}

	item := input[keys[0]]

	if len(keys) > 1 {
		nextMap, ok := item.(map[string]interface{}) 
		if ok {
			return getContent[T](nextMap, otherwise, keys[1:]...)
		}
		return otherwise
	} else {
		value, ok := item.(T) 
		if ok {
			return value
		}
		return otherwise
	}
}



func isValidPath(path string) bool {
  if path == "" {
		return false
	}

	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}
