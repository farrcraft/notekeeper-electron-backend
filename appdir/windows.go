// +build windows,!darwin

package appdir

import "os"

var systemSettingFolders = []string{os.Getenv("PROGRAMDATA")}
var globalSettingFolder = os.Getenv("APPDATA")
var cacheFolder = os.Getenv("LOCALAPPDATA")
