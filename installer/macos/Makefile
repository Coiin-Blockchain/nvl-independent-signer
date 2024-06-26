VERSION=$(shell cat ../../version)

build:
	# build the main binaries and copy them to the installer directory
	cd ../../ && make darwin
	cp ../../build/independent-signer_darwin_amd64 independent-signer_darwin_amd64

	# sign the binary
	codesign --timestamp=http://timestamp.apple.com/ts01 -s "$(COIIN_LLC)" -v independent-signer_darwin_amd64 --options=runtime

	# build the installer
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -v -o independent-signer-installer_darwin_amd64 -ldflags="-s -w -X main.Version=$(VERSION) -linkmode=external -extldflags=-mmacosx-version-min=10.14"

	# cleanup
	rm independent-signer_darwin_amd64

	# move the installer to the app directory
	cp independent-signer-installer_darwin_amd64 "Raiinmaker-Network-Validator.app/Contents/MacOS/Raiinmaker-Network-Validator"

	# cleanup
	rm independent-signer-installer_darwin_amd64

	# remove disallowed xattr
	xattr -cr "Raiinmaker-Network-Validator.app"

	# sign the installer
	codesign --timestamp=http://timestamp.apple.com/ts01 -s "$(COIIN_LLC)" -v "Raiinmaker-Network-Validator.app" --options=runtime