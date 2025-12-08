//go:build windows
package steam

import (
	"golang.org/x/sys/windows/registry"
	"log"
	"os"
)

const(
	STEAM_REGISTRY = `SOFTWARE\WOW6432Node\Valve\Steam`
	STEAM_INSTALLATION_PATH = `InstallPath`
	STEAM_LIBRARIES_VDF = `steamapps\libraryfolders.vdf`
)


func getSteamInstallDirectory() string {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, STEAM_REGISTRY, registry.QUERY_VALUE)
	if err != nil {
		log.Fatal(err)
		return "" 
	}
	defer k.Close()

	path, _, err := k.GetStringValue(STEAM_INSTALLATION_PATH)
	if err != nil {
		log.Fatal(err)
		return "" 
	}

	_, err = os.Stat(path);
	if err != nil {
		return ""
	}

	return path 
}
