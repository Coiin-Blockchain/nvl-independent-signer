package coiingui

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Coiin-Blockchain/nvl-independent-signer/installer/windows/src/coiininstaller"
	"github.com/Coiin-Blockchain/nvl-independent-signer/installer/windows/src/publickey"
)

type CoiinGUI struct {
	w              fyne.Window
	a              fyne.App
	appName        string
	coiinInstaller *coiininstaller.CoiinInstaller
}

func NewGUI(appName string, coiinInstaller *coiininstaller.CoiinInstaller) *CoiinGUI {
	a := app.New()
	w := a.NewWindow(appName)
	return &CoiinGUI{
		appName:        appName,
		a:              a,
		w:              w,
		coiinInstaller: coiinInstaller,
	}
}

func (cgui *CoiinGUI) StartInstaller() {
	cgui.w.Resize(fyne.NewSize(0, 0))

	// Check if com.coiin.independent-signer is already installed
	isInstalled, err := cgui.coiinInstaller.IsInstalled()
	if err != nil {
		log.Printf("Error encountered while checking if IndependentSigner is installed: $v\n", err)
		log.Printf("Shutting down installer in 10 seconds\n")
		time.Sleep(time.Second * 10)
		cgui.w.Close()
	}

	if isInstalled {
		installer := widget.NewLabel(fmt.Sprintf("%s is already installed.\nDo you want to uninstall or override the current version?", cgui.appName))
		cgui.w.SetContent(container.NewVBox(
			installer,
			widget.NewButton("1. Override the current version", func() {
				cgui.installGUI()
			}),
			widget.NewButton("2. Uninstall the current version", func() {
				cgui.uninstallGUI(installer)
			}),
			widget.NewButton("3. Copy your Public Key to the clipboard", func() {
				cgui.copyPublicKeyGUI(installer)
			}),
		))
	} else {
		installer := widget.NewLabel(fmt.Sprintf("Do you want to install %s?", cgui.appName))
		cgui.w.SetContent(container.NewVBox(
			installer,
			widget.NewButton("1. Yes", func() {
				cgui.installGUI()
			}),
		))
	}

	cgui.w.ShowAndRun()
}

func (cgui *CoiinGUI) showCloseButton(installer *widget.Label) {
	cgui.w.SetContent(container.NewVBox(
		installer,
		widget.NewButton("Close", func() {
			cgui.w.Close()
		}),
	))
	cgui.w.Resize(fyne.NewSize(0, 0))
}

func (cgui *CoiinGUI) uninstallGUI(installer *widget.Label) {
	defer cgui.showCloseButton(installer)
	err := cgui.coiinInstaller.Uninstall()
	if err != nil {
		installer.SetText(fmt.Sprintf("An error occurred while uninstalling %s: %s", cgui.appName, err))
		return
	}
	installer.SetText(fmt.Sprintf("%s was uninstalled successfully!", cgui.appName))
}

func (cgui *CoiinGUI) installGUI() {
	widgetLabel := widget.NewLabel(fmt.Sprintf("Installing %s ...", cgui.appName))
	cgui.w.SetContent(container.NewVBox(
		widgetLabel,
		widget.NewProgressBarInfinite(),
	))
	cgui.w.Resize(fyne.NewSize(0, 0))

	_, err := cgui.coiinInstaller.Install()
	if err != nil {
		widgetLabel = widget.NewLabel(fmt.Sprintf("An error occurred while installing %s: %s", cgui.appName, err))
		cgui.w.SetContent(container.NewVBox(
			widgetLabel,
			widget.NewButton("Close", func() {
				cgui.w.Close()
			}),
		))
		return
	}

	url := &url.URL{
		Scheme: "https",
		Host:   "coiin.io",
		Path:   "/console/verificationnodes",
	}

	publicKey, _ := publickey.CopyPublicKeyToClipboard()
	cgui.w.SetContent(container.NewVBox(
		widget.NewLabel(fmt.Sprintf("The Public Key has been copied to the clipboard:\n%s", publicKey)),
		widget.NewHyperlink("\nNavigate to the Network Validation Layer Nodes page on the Coiin Console.", url),
		widget.NewLabel("Paste the Public Key printed in the terminal window into the \"Enter Public Key\" text box and click the \"Register Node\" button."),
		widget.NewLabel(fmt.Sprintf("\n%s was installed successfully", cgui.appName)),
		widget.NewButton("Close", func() {
			cgui.w.Close()
		}),
	))
}

func (cgui *CoiinGUI) copyPublicKeyGUI(installer *widget.Label) {
	defer cgui.showCloseButton(installer)
	publickey, err := publickey.CopyPublicKeyToClipboard()
	if err != nil {
		installer.SetText(fmt.Sprintf("An error occurred while installing %s: %s", cgui.appName, err))
		return
	}
	installer.SetText(fmt.Sprintf("The Public Key has being copied to your clipboard: \n\n%s", publickey))
}
