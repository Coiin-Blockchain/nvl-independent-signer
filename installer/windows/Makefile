build:
	# build the main binaries and copy them to the installer directory
	cd ../../ && make build
	cp ../../build/independent-signer_windows_amd64.exe independent-signer_windows_amd64.exe

	# create the resource files for the installer
	go-winres make

	# build the installer
	GOOS=windows GOARCH=amd64 go build -v -o independent-signer-installer_windows_amd64.exe -ldflags="-s -w"

	# cleanup
	rm independent-signer_windows_amd64.exe
	rm rsrc_windows_386.syso
	rm rsrc_windows_amd64.syso

	# rename the installer (the installer name should end with "_windows_amd64.exe")
	mv independent-signer-installer_windows_amd64.exe independent-signer-installer.exe
