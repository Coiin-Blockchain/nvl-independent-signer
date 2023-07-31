EXECUTABLE=independent-signer
BUILD_DIR=build
VERSION=$(shell cat version)

WINDOWS=$(EXECUTABLE)_windows_amd64.exe
LINUX=$(EXECUTABLE)_linux_amd64
DARWIN=$(EXECUTABLE)_darwin_amd64
RASPBERRY=$(EXECUTABLE)_raspberry_arm


.PHONY: all clean

all: build ## Build and run tests

build: windows linux darwin raspberry ## Build binaries
	@echo version: $(VERSION)

windows: $(WINDOWS) ## Build for Windows

linux: $(LINUX) ## Build for Linux

darwin: $(DARWIN) ## Build for Darwin (macOS)

raspberry: $(RASPBERRY) ## Build for Raspberry (Linux)

$(WINDOWS):
	env GOOS=windows GOARCH=amd64 go build -v -o $(BUILD_DIR)/$(WINDOWS) -ldflags="-s -w -X main.Version=$(VERSION)"

$(LINUX):
	env GOOS=linux GOARCH=amd64 go build -v -o $(BUILD_DIR)/$(LINUX) -ldflags="-s -w -X main.Version=$(VERSION)"

$(DARWIN):
	env GOOS=darwin GOARCH=amd64 go build -v -o $(BUILD_DIR)/$(DARWIN) -ldflags="-s -w -X main.Version=$(VERSION)"

$(RASPBERRY):
	env GOOS=linux GOARCH=arm GOARM=7 go build -v -o $(BUILD_DIR)/$(RASPBERRY) -ldflags="-s -w -X main.Version=$(VERSION)"

clean: ## Remove previous build
	rm -f $(BUILD_DIR)/$(WINDOWS) $(BUILD_DIR)/$(LINUX) $(BUILD_DIR)/$(DARWIN) $(BUILD_DIR)/$(RASPBERRY)

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
