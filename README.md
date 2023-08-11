# Getting Started

## Requirements
* Create a [Coiin Console](https://coiin.io/console) account
* Download the newest release of [nvl-independent-signer](https://github.com/Coiin-Blockchain/nvl-independent-signer/releases) for your OS

## Register Independent Signer

1. From the terminal, navigate to the folder where you downloaded, and run the independent signer script. The first time will generate a new signing key. Once the new signing key is generated, the Public Key will be printed to the terminal. Copy the Public Key.
    - Note: Some operating systems like macOS may require updating default security settings and granting execute permission to execute the `independent-signer` script.

2. Navigate to the [Register an Independent Node](https://coiin.io/console/verificationnodes) section on the Network Validation Layer Nodes page of the Coiin Console. Paste the Public Key printed in the terminal window from Step 1 into the "Enter Public Key" text box and click the "Register Node" button.

3. After the signing key is generated and the public key is saved to your Coiin Console account, you can run `independent-signer` from the terminal at any time to sign the latest NVL Proxy block and post it to the NVL Proxy.
    - Note: it is recommended to set up a cron job to execute the `independent-signer` script once every 30 minutes.

Your Public and Signing Keys are saved in:

    # Windows
    C:\Users\<Username>\AppData\Roaming\coiin\nvl\independent-signer
    
    # Linux/Mac
    ~/.config/coiin/nvl/independent-signer

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
