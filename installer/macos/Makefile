build:
	# build the main binaries and copy them to the installer directory
	cd ../../ && make build
	cp ../../build/independent-signer_darwin_amd64 independent-signer_darwin_amd64

	# sign the binary
	codesign --timestamp=http://timestamp.apple.com/ts01 -s "$(COIIN_LLC)" -v independent-signer_darwin_amd64 --options=runtime

	# build the installer
	GOOS=darwin go build -v -o independent-signer-installer_darwin_amd64 -ldflags="-s -w"

	# cleanup
	rm independent-signer_darwin_amd64

	# move the installer to the app directory
	cp independent-signer-installer_darwin_amd64 independent-signer_darwin_amd64.app/Contents/MacOS/independent-signer_darwin_amd64

	# cleanup
	rm independent-signer-installer_darwin_amd64

	# sign the installer
	codesign --timestamp=http://timestamp.apple.com/ts01 -s "$(COIIN_LLC)" -v independent-signer_darwin_amd64.app --options=runtime