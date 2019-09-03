// Package appdir based on: https://github.com/shibukawa/configdir
package appdir

import "path/filepath"

// AppDataPath returns the path where application data is stored
func AppDataPath() string {
	basePath := filepath.Clean(globalSettingFolder)
	path := filepath.Join(basePath, AppName)
	return path
}
