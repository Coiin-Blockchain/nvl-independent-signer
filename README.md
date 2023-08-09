# Getting Started

## Requirements
* Create Coiin [console](https://coiin.io/console) account
* [Register](https://coiin.io/console/verificationnodes) with NVL Proxy

Download the newest release of [nvl-independent-signer](https://github.com/Coiin-Blockchain/nvl-independent-signer/releases)

## Register Independent Signer

Run the independent signer from the terminal. The first time will generate a new signing key and ask
for a registration ID. Once the new signing key is generated, the public key will be printed to the terminal.

Go to the [independent registration](https://coiin.io/console/validationnodes) page on the coiin console. Click "BUTTON NAME" to begin the process. 
When prompted, paste the public key printed in console into the text box and click "Register". Once
registered, the console will give you a registration ID. Paste that into your terminal when prompted.

After the signing key is generated and registration ID is saved, you can run `independent-signer` at any
time to sign the latest NVL Proxy block and post it to the NVL Proxy.

It is recommended to set up a cron job to execute the `independent-signer` once every 30 minutes

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
