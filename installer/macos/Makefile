build:
	cd ../../ && make build
	cp ../../build/independent-signer_darwin_amd64 independent-signer_darwin_amd64

	# -X main.VersionInfo.FileDescription='Coiin Network Indep. Validator' -X main.Copyright='Copyright © 2023 Huh Coiin LLC'
	GOOS=darwin GOARCH=amd64 go build -v -o independent-signer_macos -ldflags="-s -w"