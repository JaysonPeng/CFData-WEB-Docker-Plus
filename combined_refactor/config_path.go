package main

import (
	"os"
	"path/filepath"
	"strings"
)

const configFileName = "cfdata-config.json"

// getConfigFilePath returns the path used by the Web server for
// cfdata-config.json.
//
// CFDATA_CONFIG_PATH can be used to explicitly place the configuration file,
// which is useful for Docker/NAS persistent storage. When the environment
// variable is not set, the original behavior is preserved: the file is stored
// beside the executable.
func getConfigFilePath() string {
	if path := strings.TrimSpace(os.Getenv("CFDATA_CONFIG_PATH")); path != "" {
		return path
	}

	exe, err := os.Executable()
	if err == nil && exe != "" {
		return filepath.Join(filepath.Dir(exe), configFileName)
	}

	return configFileName
}
