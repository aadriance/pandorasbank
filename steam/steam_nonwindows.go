//go:build !windows
package steam

const (
	STEAM_LIBRARIES_VDF = `steamapps/libraryfolders.vdf`
)

func getSteamInstallDirectory() string {
	//TODO ?
	return ""
}
