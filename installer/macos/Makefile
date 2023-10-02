build:
	# build the main binaries and copy them to the installer directory
	cd ../../ && make build
	cp ../../build/independent-signer_darwin_amd64 independent-signer_darwin_amd64

	# build the installer
	GOOS=darwin GOARCH=amd64 go build -v -o independent-signer-installer_darwin_amd64 -ldflags="-s -w"

	# cleanup
	rm independent-signer_darwin_amd64

	# move the installer to the app directory
	mv independent-signer-installer_darwin_amd64 independent-signer_darwin_amd64