#!/bin/bash

# Define the current path
current_path=$(pwd)

# Define the full path to the executable
executable_path="$current_path/independent-signer_darwin_amd64"

# Define the name of the scheduled task
task_name="IndependentSigner"

# Create a launchd plist to run the program every 30 minutes
cat << EOF > "$task_name.plist"
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>$task_name</string>
    <key>ProgramArguments</key>
    <array>
        <string>$executable_path</string>
    </array>
    <key>StartInterval</key>
    <integer>1800</integer>
    <key>RunAtLoad</key>
    <true/>
</dict>
</plist>
EOF

# Tip:
# If something goes wrong, you can try this: launchctl debug IndependentSigner --stderr error.log && launchctl kickstart -k IndependentSigner
# This will run the program in debug mode and save the output to error.log
# Load the launchd plist
launchctl remove "$task_name"
launchctl load -F "$task_name.plist"

# Run the program and save its output to a temporary file
echo ""
"$executable_path" > temp.txt 2>&1

# Find the line containing "Public Key:" and extract the fixed-size value
grep -o "Public Key: [a-zA-Z0-9]\{130\}" temp.txt > extracted.txt

# Extract the Public Key value from the line
publickey=$(head -n 1 extracted.txt | cut -c 13-)

# Copy the extracted Public Key value to the clipboard
echo "$publickey" | pbcopy

# Display a message indicating that the Public Key has been copied to the clipboard
echo "The Public Key has been copied to the clipboard:"
echo "$publickey"

# Display instructions for the user
echo ""
echo "Navigate to https://coiin.io/console/verificationnodes on the Network Validation Layer Nodes page of the Coiin Console."
echo "Paste the Public Key printed in the terminal window into the \"Enter Public Key\" text box and click the \"Register Node\" button."
echo ""

# Clean up temporary files (optional)
rm temp.txt extracted.txt "$task_name.plist"