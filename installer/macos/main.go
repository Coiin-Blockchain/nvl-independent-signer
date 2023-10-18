package main

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

//go:embed independent-signer_darwin_amd64
var independentSigner embed.FS

//go:embed script.sh
var script embed.FS

func uninstall() error {
	cmd := exec.Command("launchctl", "remove", "IndependentSigner")
	return cmd.Run()
}

func install() (string, error) {
	// Get the default path for the independent-signer_windows_amd64.exe
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir, err = os.Getwd()
		if err != nil {
			return "", fmt.Errorf("could not find working directory: %w", err)
		}
	}
	defaultPath := filepath.Join(configDir, "coiin", "nvl", "independent-signer")

	// Set the paths for the temporary files
	tempExe := filepath.Join(defaultPath, "independent-signer_darwin_amd64")
	err = os.MkdirAll(filepath.Dir(tempExe), 0755)
	if err != nil {
		return "", err
	}

	tempBat := filepath.Join(defaultPath, "script.sh")
	err = os.MkdirAll(filepath.Dir(tempBat), 0755)
	if err != nil {
		return "", err
	}

	// Read the independent-signer_darwin_amd64 from the embed.FS
	exeContent, err := independentSigner.ReadFile("independent-signer_darwin_amd64")
	if err != nil {
		return "", err
	}

	// Read the script.sh from the embed.FS
	batContent, err := script.ReadFile("script.sh")
	if err != nil {
		return "", err
	}

	// Write the independent-signer_darwin_amd64 to the temp directory
	err = os.WriteFile(tempExe, exeContent, 0755)
	if err != nil {
		return "", err
	}

	// Write the script.sh to the temp directory
	err = os.WriteFile(tempBat, batContent, 0755)
	if err != nil {
		return "", err
	}

	// Change the working directory to defaultPath
	err = os.Chdir(defaultPath)
	if err != nil {
		return "", err
	}

	// Run the script.sh
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

func uninstallGUI(installer *widget.Label) {
	defer showCloseButton(installer)
	err := uninstall()
	if err != nil {
		installer.SetText(fmt.Sprintf("An error occurred while uninstalling Independent Signer: %s", err))
		return
	}
	installer.SetText("Independent Signer was uninstalled successfully!")
}

func installGUI(installer *widget.Label) {
	defer showCloseButton(installer)
	installer.SetText("Installing Independent Signer ...")
	msg, err := install()
	if err != nil {
		installer.SetText(fmt.Sprintf("An error occurred while installing Independent Signer: %s", err))
		return
	}
	installer.SetText(fmt.Sprintf("%s\n\nIndependent Signer was installed successfully!", msg))
}

func showCloseButton(installer *widget.Label) {
	w.SetContent(container.NewVBox(
		installer,
		widget.NewButton("Close", func() {
			os.Exit(1)
		}),
	))
}

var (
	a = app.New()
	w = a.NewWindow("independent-signer-installer")
)

func main() {
	w.Resize(fyne.NewSize(250, 300))

	// Check if IndependentSigner is already installed
	cmd := exec.Command("launchctl", "list", "IndependentSigner")
	output, err := cmd.Output()
	if err == nil {
		if len(output) > 0 {

			installer := widget.NewLabel("IndependentSigner is already installed. Do you want to uninstall or override the current version?")
			w.SetContent(container.NewVBox(
				installer,
				widget.NewButton("1. Override the current version", func() {
					installGUI(installer)
				}),
				widget.NewButton("2. Uninstall the current version", func() {
					uninstallGUI(installer)
				}),
			))

		}

	} else {
		installer := widget.NewLabel("Do you want to install the IndependentSigner?")
		w.SetContent(container.NewVBox(
			installer,
			widget.NewButton("1. Yes", func() {
				installGUI(installer)
			}),
		))

	}

	w.ShowAndRun()

}
