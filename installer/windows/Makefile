build:
	cd ../../ && make build
	cp ../../build/independent-signer_windows_amd64.exe independent-signer_windows_amd64.exe

	# -X main.VersionInfo.FileDescription='Coiin Network Indep. Validator' -X main.Copyright='Copyright © 2023 Huh Coiin LLC'
	GOOS=windows GOARCH=amd64 go build -v -o independent-signer_windows.exe -ldflags="-s -w"