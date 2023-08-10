# Getting Started

## Requirements
* Create [Coiin Console](https://coiin.io/console) account
* [Register](https://coiin.io/console/verificationnodes) with NVL Proxy

Download the newest release of [nvl-independent-signer](https://github.com/Coiin-Blockchain/nvl-independent-signer/releases)

## Register Independent Signer

1. Run the independent signer from the terminal. The first time will generate a new signing key. Once the new signing key is generated, the public key will be printed to the terminal.
2. Go to the [Register an Independent Node](https://coiin.io/console/verificationnodes) section on the Network Validation Layer Nodes page of the Coiin Console. Paste the public key printed in the terminal window into the Enter Public Key text box and click the "Register Node" button.
3. After the signing key is generated and the public key is saved to your Coiin Console account, you can run `independent-signer` at any time to sign the latest NVL Proxy block and post it to the NVL Proxy.

Note: It is recommended to set up a cron job to execute the `independent-signer` once every 30 minutes.

Keys and registration ID are saved in

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
