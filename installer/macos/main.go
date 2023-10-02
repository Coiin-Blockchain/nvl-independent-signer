package main

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

//go:embed independent-signer_darwin_amd64
var independentSigner embed.FS

//go:embed script.sh
var script embed.FS

func main() {
	// Clear the console
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	// Check if IndependentSigner is already installed
	cmd = exec.Command("launchctl", "list", "IndependentSigner")
	output, err := cmd.Output()
	if err == nil {
		if len(output) > 0 {
			// IndependentSigner is already installed, prompt the user to uninstall or override
			fmt.Print("IndependentSigner is already installed. Do you want to uninstall or override the current version?\n\n")
			fmt.Println("1. Override the current version")
			fmt.Println("2. Uninstall the current version")
			reader := bufio.NewReader(os.Stdin)
			var answer int
			// Loop until the user enters a valid input
			for {
				fmt.Print("Enter your choice (1 or 2): ")
				input, err := reader.ReadString('\n')
				if err != nil {
					log.Fatal(err)
				}

				// Convert the input to an integer
				input = strings.TrimSpace(input)
				answer, err = strconv.Atoi(input)

				// Check if the input is valid
				if err == nil && (answer == 1 || answer == 2) {
					break
				}
				fmt.Println("Invalid input. Please enter 1 or 2.")
			}
			// Override the current version
			if answer == 2 {
				cmd = exec.Command("launchctl", "remove", "IndependentSigner")
				cmd.Run()
				fmt.Print("\nIndependent Signer was uninstalled successfully!\n")
				os.Exit(0)
			}
		}
	}

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
	cmd = exec.Command(tempBat)
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
}
