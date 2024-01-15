package main

import (
	"bytes"
	"embed"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/atotto/clipboard"
)

//go:embed independent-signer_darwin_amd64
var independentSigner embed.FS

//go:embed script.sh
var script embed.FS

func getDefaultPath() (string, error) {
	// Get the default path for the independent-signer_windows_amd64.exe
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir, err = os.Getwd()
		if err != nil {
			return "", fmt.Errorf("could not find working directory: %w", err)
		}
	}
	return filepath.Join(configDir, "coiin", "nvl", "independent-signer"), nil
}

func uninstall() error {
	cmd := exec.Command("launchctl", "remove", "com.coiin.independent-signer")
	return cmd.Run()
}

func install() (string, error) {

	defaultPath, err := getDefaultPath()
	if err != nil {
		return "", err
	}

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

func showCloseButton(installer *widget.Label) {
	w.SetContent(container.NewVBox(
		installer,
		widget.NewButton("Close", func() {
			os.Exit(1)
		}),
	))
	w.Resize(fyne.NewSize(0, 0))
}

func copyPublicKeyToClipboard() (string, error) {
	defaultPath, err := getDefaultPath()
	if err != nil {
		return "", err
	}
	publicKeyFile := filepath.Join(defaultPath, "public-key")

	// Read the public key from the default path
	content, err := os.ReadFile(publicKeyFile)
	if err != nil {
		return "", err
	}
	content = bytes.TrimSpace(content)

	err = clipboard.WriteAll(string(content))
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func uninstallGUI(installer *widget.Label) {
	defer showCloseButton(installer)
	err := uninstall()
	if err != nil {
		installer.SetText(fmt.Sprintf("An error occurred while uninstalling %s: %s", appName, err))
		return
	}
	installer.SetText(fmt.Sprintf("%s was uninstalled successfully!", appName))
}

func installGUI(isUpgrade bool) {
	gui := widget.NewLabel(fmt.Sprintf("Installing %s ...", appName))
	w.SetContent(container.NewVBox(
		gui,
		widget.NewProgressBarInfinite(),
	))
	w.Resize(fyne.NewSize(0, 0))

	_, err := install()
	if err != nil {
		gui = widget.NewLabel(fmt.Sprintf("An error occurred while installing %s: %s", appName, err))
		w.SetContent(container.NewVBox(
			gui,
			widget.NewButton("Close", func() {
				os.Exit(1)
			}),
		))
		return
	}

	url := &url.URL{
		Scheme: "https",
		Host:   "coiin.io",
		Path:   "/console/verificationnodes",
	}

	publicKey, _ := copyPublicKeyToClipboard()

	if isUpgrade {
		w.SetContent(container.NewVBox(
			widget.NewLabel(fmt.Sprintf("The Public Key has been copied to the clipboard:\n%s", publicKey)),
			widget.NewLabel(fmt.Sprintf("\n%s was upgraded successfully. No further action is necessary", appName)),
			widget.NewButton("Close", func() {
				os.Exit(1)
			}),
		))
		return
	}

	w.SetContent(container.NewVBox(
		widget.NewLabel(fmt.Sprintf("The Public Key has been copied to the clipboard:\n%s", publicKey)),
		widget.NewHyperlink("\nNavigate to the Network Validation Layer Nodes page on the Coiin Console.", url),
		widget.NewLabel("Paste the Public Key printed in the terminal window into the \"Enter Public Key\" text box and click the \"Register Node\" button."),
		widget.NewLabel(fmt.Sprintf("\n%s was installed successfully", appName)),
		widget.NewButton("Close", func() {
			os.Exit(1)
		}),
	))
}

func copyPublicKeyGUI(installer *widget.Label) {
	defer showCloseButton(installer)
	publickey, err := copyPublicKeyToClipboard()
	if err != nil {
		installer.SetText(fmt.Sprintf("An error occurred while installing %s: %s", appName, err))
		return
	}
	installer.SetText(fmt.Sprintf("The Public Key has being copied to your clipboard: \n\n%s", publickey))
}

const (
	appName = "Coiin Network Validator"
)

var (
	Version = "v0.0.0"
	a       = app.New()
	w       = a.NewWindow(appName)
)

func main() {
	w.Resize(fyne.NewSize(0, 0))

	// Check if com.coiin.independent-signer is already installed
	cmd := exec.Command("launchctl", "list", "com.coiin.independent-signer")
	output, err := cmd.Output()
	if err == nil {
		if len(output) > 0 {

			installer := widget.NewLabel(fmt.Sprintf("%s is already installed.\nDo you want to uninstall or override the current version?", appName))
			w.SetContent(container.NewVBox(
				installer,
				widget.NewButton(fmt.Sprintf("1. Override the current version with %s", Version), func() {
					installGUI(true)
				}),
				widget.NewButton("2. Uninstall the current version", func() {
					uninstallGUI(installer)
				}),
				widget.NewButton("3. Copy your Public Key to the clipboard", func() {
					copyPublicKeyGUI(installer)
				}),
			))

		}

	} else {
		installer := widget.NewLabel(fmt.Sprintf("Do you want to install the %s - %s?", appName, Version))
		w.SetContent(container.NewVBox(
			installer,
			widget.NewButton("1. Yes", func() {
				installGUI(false)
			}),
		))

	}

	w.ShowAndRun()

}
