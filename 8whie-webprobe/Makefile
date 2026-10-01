# 8WHIE WebProbe Makefile
# Project: 8WHIE WebProbe
# Author: 8WHIE (Android • Termux • Linux • Security Research)

BINARY_NAME=webprobe
CMD_PATH=./cmd/webprobe
BUILD_DIR=./bin
VERSION=1.0.0

.PHONY: all build test vet fmt clean install cross-compile

all: fmt vet test build

build:
	@echo "[*] Building $(BINARY_NAME) v$(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)
	@echo "[+] Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

test:
	@echo "[*] Running test suite..."
	go test -v -race ./tests/...

vet:
	@echo "[*] Running go vet..."
	go vet ./...

fmt:
	@echo "[*] Running gofmt..."
	gofmt -w -s .

clean:
	@echo "[*] Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR) results.json results.csv results.txt

install: build
	@echo "[*] Installing to GOPATH/bin..."
	go install $(CMD_PATH)

cross-compile:
	@echo "[*] Building cross-platform distribution binaries..."
	@mkdir -p $(BUILD_DIR)/release
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/release/$(BINARY_NAME)-linux-amd64 $(CMD_PATH)
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/release/$(BINARY_NAME)-linux-arm64 $(CMD_PATH)
	GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-s -w" -o $(BUILD_DIR)/release/$(BINARY_NAME)-linux-armv7-termux $(CMD_PATH)
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/release/$(BINARY_NAME)-darwin-amd64 $(CMD_PATH)
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/release/$(BINARY_NAME)-darwin-arm64 $(CMD_PATH)
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/release/$(BINARY_NAME)-windows-amd64.exe $(CMD_PATH)
	@echo "[+] Cross-compilation complete in $(BUILD_DIR)/release"
