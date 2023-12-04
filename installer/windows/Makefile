build:
	# build the main binaries and copy them to the installer directory
	cd ../../ && make windows
	cp ../../build/independent-signer_windows_amd64.exe independent-signer_windows_amd64.exe

	# create the resource files for the installer
	go-winres make

	# https://stackoverflow.com/questions/74334529/go-fyne-project-cant-cross-compile-from-linux-to-windows
	# https://developer.fyne.io/started/cross-compiling
	# build the installer
	env CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=/usr/bin/x86_64-w64-mingw32-gcc go build -v -o independent-signer-installer_windows_amd64.exe -ldflags="-s -w -H=windowsgui"

	# cleanup
	rm independent-signer_windows_amd64.exe
	rm rsrc_windows_386.syso
	rm rsrc_windows_amd64.syso

	# rename the installer (the installer name should end with "_windows_amd64.exe")
	mv independent-signer-installer_windows_amd64.exe "Coiin Network Validator.exe"