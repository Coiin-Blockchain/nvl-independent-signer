package main

import (
	"embed"

	"github.com/Coiin-Blockchain/nvl-independent-signer/installer/windows/src/coiingui"
	"github.com/Coiin-Blockchain/nvl-independent-signer/installer/windows/src/coiininstaller"
)

//go:embed independent-signer_windows_amd64.exe
var independentSigner embed.FS

//go:embed script.bat
var script embed.FS

const appName = "Coiin Network Validator"

func main() {
	coiinInstaller := coiininstaller.NewCoiinInstaller(&independentSigner, &script)
	coiinGui := coiingui.NewGUI(appName, coiinInstaller)
	coiinGui.StartInstaller()
}
