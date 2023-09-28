package main

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// This project should be able to embed the independent-signer_windows_amd64.exe while building, and extract it when it's running, than run a .bat file that will start the independent-signer_windows_amd64.exe

//go:embed independent-signer_windows_amd64.exe
var independentSigner embed.FS

//go:embed script.bat
var script embed.FS

func main() {
	// Default path for the independent-signer_windows_amd64.exe
	defaultPath := `C:\Users\Public\Independent-Signer`

	// Create the directory if it does not exist
	err := os.MkdirAll(defaultPath, 0755)
	if err != nil {
		log.Fatal(err)
	}

	// Set the paths for the temporary files
	tempExe := filepath.Join(defaultPath, "independent-signer_windows_amd64.exe")
	tempBat := filepath.Join(defaultPath, "script.bat")

	// Read the independent-signer_windows_amd64.exe from the embed.FS
	exeContent, err := independentSigner.ReadFile("independent-signer_windows_amd64.exe")
	if err != nil {
		log.Fatal(err)
	}

	// Read the script.bat from the embed.FS
	batContent, err := script.ReadFile("script.bat")
	if err != nil {
		log.Fatal(err)
	}

	// Write the independent-signer_windows_amd64.exe to the temp directory
	err = os.WriteFile(tempExe, exeContent, 0755)
	if err != nil {
		log.Fatal(err)
	}

	// Write the script.bat to the temp directory
	err = os.WriteFile(tempBat, batContent, 0755)
	if err != nil {
		log.Fatal(err)
	}

	// Print a message to the console
	fmt.Print("Installing Independent Signer ...\n\n")

	// Wait for 2 seconds
	time.Sleep(2 * time.Second)

	// Change the working directory to defaultPath
	err = os.Chdir(defaultPath)
	if err != nil {
		log.Fatal(err)
	}

	// Run the script.bat
	cmd := exec.Command(tempBat)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Print the output of the script to the console
	fmt.Println(out.String())

	// Print a message to the console
	fmt.Print("Independent Signer installed successfully!\n\n")

	// Pause the program
	var input string
	fmt.Scanln(&input)
}
