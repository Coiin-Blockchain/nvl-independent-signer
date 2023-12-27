package coiinpath

import (
	"fmt"
	"os"
	"path/filepath"

	coiincommon "github.com/Coiin-Blockchain/nvl-independent-signer/installer/windows/src/coincommon"
)

func GetDefaultPath() (string, error) {
	configDir, err := GetUserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, coiincommon.CoiinFolderName, "nvl", "independent-signer"), nil
}

func GetUserConfigDir() (string, error) {
	// Get the default path for the independent-signer_windows_amd64.exe
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir, err = os.Getwd()
		if err != nil {
			return "", fmt.Errorf("could not find working directory: %w", err)
		}
	}
	return configDir, nil
}
