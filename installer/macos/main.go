package main

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// This project should be able to embed the independent-signer_darwin_amd64 while building, and extract it when it's running, than run a .sh file that will start the independent-signer_darwin_amd64

//go:embed independent-signer_darwin_amd64
var independentSigner embed.FS

//go:embed script.sh
var script embed.FS

func main() {
	// Get the default path for the independent-signer_windows_amd64.exe
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir, err = os.Getwd()
		if err != nil {
			log.Fatal("could not find working directory")
		}
	}
	defaultPath := filepath.Join(configDir, "coiin", "nvl", "independent-signer")

	// Set the paths for the temporary files
	tempExe := filepath.Join(defaultPath, "independent-signer_darwin_amd64")
	tempBat := filepath.Join(defaultPath, "script.sh")

	// Read the independent-signer_darwin_amd64 from the embed.FS
	exeContent, err := independentSigner.ReadFile("independent-signer_darwin_amd64")
	if err != nil {
		log.Fatal(err)
	}

	// Read the script.sh from the embed.FS
	batContent, err := script.ReadFile("script.sh")
	if err != nil {
		log.Fatal(err)
	}

	// Write the independent-signer_darwin_amd64 to the temp directory
	err = os.WriteFile(tempExe, exeContent, 0755)
	if err != nil {
		log.Fatal(err)
	}

	// Write the script.sh to the temp directory
	err = os.WriteFile(tempBat, batContent, 0755)
	if err != nil {
		log.Fatal(err)
	}

	// Print a message to the console
	fmt.Print("Installing Independent Signer ...\n\n")

	// Change the working directory to defaultPath
	err = os.Chdir(defaultPath)
	if err != nil {
		log.Fatal(err)
	}

	// Run the script.sh
	cmd := exec.Command(tempBat)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	filePath := fmt.Sprintf("%s/%s", defaultPath, "instructions")
	err = os.WriteFile(filePath, []byte(out.String()), 0644)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}
	cmd = exec.Command("open", filePath)
	cmd.Run()

	// Print the output of the script to the console
	fmt.Println(out.String())

	// Print a message to the console
	fmt.Print("Independent Signer installed successfully!\n\n")

	// Pause the program
	var input string
	fmt.Scanln(&input)
}
