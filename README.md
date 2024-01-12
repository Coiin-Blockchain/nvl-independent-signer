# Getting Started

## Requirements
* Create a [Coiin Console](https://coiin.io/console) account

## Register Independent Signer From an Installer (Easy Install)

1. Download the latest release of the [Coiin Network Validator](https://github.com/Coiin-Blockchain/nvl-independent-signer/releases) installer for your OS.
   - Note: macOS users should download the file ending in .dmg, while Windows users should download the file ending in .exe
  
2. Run the installer. 
    - If you don't see your Public Key after running the installer, try running it again.
    
    <img src = "/assets/installer.png" width=50% height=50%>
    <img src = "/assets/installcomplete.png" width=50% height=50%>
   
3. Paste the Public Key from the installer into the [Validation Nodes](https://coiin.io/console/verificationnodes) page of the Coiin Console and click Register Node. 
   - This will associate your Public Key with your Coiin Console account to ensure you get rewarded for mining NVL blocks.
4. You're done! The independent signer script will automatically run every 30 minutes and look for new NVL blocks to sign with your Public Key

## Register Independent Signer From a Command Line Interface

* Download the latest release of the [Independent-Signer](https://github.com/Coiin-Blockchain/nvl-independent-signer/releases) script for your OS, then follow the instructions below.

<details>

<summary>Windows</summary>

### Windows Installation Instructions

1. From the terminal, navigate to the folder where you downloaded the file and execute the independent signer script. The first time will generate a new signing key. Once the new signing key is generated, the Public Key will be printed to the terminal. Copy the Public Key.

2. Navigate to the [Register an Independent Node](https://coiin.io/console/verificationnodes) section on the Network Validation Layer Nodes page of the Coiin Console. Paste the Public Key printed in the terminal window from Step 1 into the "Enter Public Key" text box and click the "Register Node" button.

3. After the signing key is generated and the public key is saved to your Coiin Console account, you can run `independent-signer` from the terminal at any time to sign the latest NVL Proxy block and post it to the NVL Proxy.
    - Note: to make this easy, it is recommended to set up a cron job to execute the `independent-signer` script once every 30 minutes.

Your Public and Signing Keys are saved in:

    C:\Users\<Username>\AppData\Roaming\coiin\nvl\independent-signer

</details>

<details>

<summary>Linux and macOS</summary>

### Linux and macOS Installation Instructions

1. From a web browser, download the Independent Signer script by navigating to: https://github.com/Coiin-Blockchain/nvl-independent-signer/releases
    - Select the script for your computer, e.g. _independent-signer_darwin_amd64_ for macOS users or __independent-signer_linux_amd64_ for Linux users
    - Save this to a preferred location where it won’t be deleted.

2. Open a new Terminal window (Terminal can be found in the Applications > Utilities folder)
    - Navigate to the directory where you saved the script by typing
      
        ```
        cd [the_filepath_you_saved_the_script_to]
        ```
        
    - Or if you're unsure of where to locate the file path you can simply drag and drop the file onto the Terminal and the file path will be shown.
    - Execute the script by typing the filename and pressing return, e.g.

            independent-signer_darwin_amd64

        - Note: if you receive an error, you may need to re-permission the script as an executable file by typing:

            ```
            chmod +x independent-signer_darwin_amd64 
            ```
            and then continue by re-executing step 2a. 
        - Or on macOS, you may need to allow the file to be opened by selecting the Apple menu  > System Settings, then click Privacy & Security in the sidebar. (You may need to scroll down.)
          Open Privacy & Security settings. Go to Security, click the pop-up menu next to “Allow applications downloaded from,” then choose the sources from which you’ll allow software to be installed:
    - The script will run, generating a Public Key, and will attempt to sign an NVL block, but will fail - that’s ok! You’ll fix that in just a moment by registering your Public Key to your Coiin Console account. For now, just copy the Public Key generated by the script.
    
3. From a web browser, log into the Coiin Console by navigating to: https://coiin.io/console/
    - Navigate to the Validation Nodes page from the menu
    - Under Independent Node Status, paste the Public Key value from your Terminal window into the field “Enter Public Key (generated from NVL script)”
    - Click the Register Node button, accept the Terms of Use, and you should see the Node Identity status update to Registered.


4. You’re almost done! Now you just need to run a cron job so that the Independent Signer script automatically runs every 30 minutes and signs each new block created by the NVL. Navigate back to your Terminal window from earlier.

    - Type this command:
        ```
        crontab -e
        ```
    - This will open the cron editor. Now type  (shift i), to enter insert mode or vim editor where you can then type the command:
        ```
        */30 * * * * [the_filepath_you_saved_the_script_to]
        ```
        Or if you're unsure of where to locate the file path you can simply drag and drop the file onto the Terminal and the file path will be shown. Make sure you use a space between the last asterisk and file name when you drop the file in.

    - If you're not in command mode (where you can type commands directly into vim), press the Esc key on your keyboard. This ensures you're in command mode.
Type the following and press return to save and close the cron editor
        ```
        :wq
        ```
        (The :wq command is a combination of two commands: :w (which saves the changes) and :q (which quits the editor).


5. That’s it! Now your computer will run the Independent Signer script every 30 minutes, fetching and signing the most recent NVL block from the Coiin blockchain. You can see the latest block that was signed on your coiin console account.

Note: if you lose access to your Coiin Console account and need to reset your password, you will also need to re-register your independent node by first clearing the contents of ~/.config/coiin/nvl/independent-signer
and return to Step 2 of the process above to create a new Public Key for your Independent node.

To verify if your crontab is functioning, type this command:

```
crontab -l
```

It should read something similar to: 

```
*30/* * * * /Users/~/[the_filepath_you_saved_the_script_to]/independent-signer_darwin_amd64
```

Your Public and Signing Keys are saved in:

    # Linux/Mac
    ~/.config/coiin/nvl/independent-signer

</details>

# Build from source

## Requirements
* [Git](https://git-scm.com/) or [GitHub client](https://desktop.github.com/)
* [Go v1.20+](https://go.dev/dl/)


Clone source code

    mkdir ~/go/src/github.com/Coiin-Blockchain
    cd ~/go/src/github.com/Coiin-Blockchain/nvl-independent-signer
    git clone https://github.com/Coiin-Blockchain/nvl-independent-signer.git

Build

    go mod tidy
    go build

# Support

* [Submit issue](https://github.com/Coiin-Blockchain/nvl-independent-signer/issues)
