package publickey

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/Coiin-Blockchain/nvl-independent-signer/installer/windows/src/coiinpath"
	"github.com/atotto/clipboard"
)

func CopyPublicKeyToClipboard() (string, error) {
	defaultPath, err := coiinpath.GetDefaultPath()
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
