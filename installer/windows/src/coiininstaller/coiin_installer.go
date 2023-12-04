package coiininstaller

import (
	"bytes"
	"embed"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Coiin-Blockchain/nvl-independent-signer/installer/windows/src/coiinpath"
	coiincommon "github.com/Coiin-Blockchain/nvl-independent-signer/installer/windows/src/coincommon"
)

type CoiinInstaller struct {
	independentSigner, script *embed.FS
}

func NewCoiinInstaller(independentSigner, script *embed.FS) *CoiinInstaller {
	return &CoiinInstaller{
		independentSigner: independentSigner,
		script:            script,
	}
}

func (ci *CoiinInstaller) IsInstalled() (bool, error) {
	cmd := exec.Command("powershell", "-NoProfile", "-WindowStyle", "hidden", "Get-ScheduledTask", "|", "?", "TaskName", "-eq", "IndependentSigner")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	return len(output) > 0, nil
}

func (ci *CoiinInstaller) Install() (string, error) {
	defaultPath, err := coiinpath.GetDefaultPath()
	if err != nil {
		return "", err
	}

	// Set the paths for the temporary files
	tempExe := filepath.Join(defaultPath, "independent-signer_windows_amd64.exe")
	err = os.MkdirAll(filepath.Dir(tempExe), 0755)
	if err != nil {
		return "", err
	}

	tempBat := filepath.Join(defaultPath, "script.bat")
	err = os.MkdirAll(filepath.Dir(tempBat), 0755)
	if err != nil {
		return "", err
	}

	// Read the independent-signer_windows_amd64.exe from the embed.FS
	exeContent, err := ci.independentSigner.ReadFile("independent-signer_windows_amd64.exe")
	if err != nil {
		return "", err
	}

	// Read the script.bat from the embed.FS
	batContent, err := ci.script.ReadFile("script.bat")
	if err != nil {
		return "", err
	}

	// Write the independent-signer_windows_amd64.exe to the temp directory
	err = os.WriteFile(tempExe, exeContent, 0755)
	if err != nil {
		return "", err
	}

	// Write the script.bat to the temp directory
	err = os.WriteFile(tempBat, batContent, 0755)
	if err != nil {
		return "", err
	}

	// Wait for 2 seconds
	time.Sleep(2 * time.Second)

	// Change the working directory to defaultPath
	err = os.Chdir(defaultPath)
	if err != nil {
		return "", err
	}

	// Run the script.bat
	cmd := exec.Command(tempBat)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Run()
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

func (ci *CoiinInstaller) Uninstall() error {
	err := removeInstalledFiles()
	if err != nil {
		return err
	}

	cmd := exec.Command("powershell", "-NoProfile", "-WindowStyle", "hidden", "Unregister-ScheduledTask", "-TaskName", "IndependentSigner", "-Confirm:$false")
	cmd.Stdout = nil
	return cmd.Run()
}

func removeInstalledFiles() error {
	configDir, err := coiinpath.GetUserConfigDir()
	if err != nil {
		return err
	}

	coiinAppPath := filepath.Join(configDir, coiincommon.CoiinFolderName)
	err = os.RemoveAll(coiinAppPath)
	if err != nil {
		return err
	}

	return nil
}
