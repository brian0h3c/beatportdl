BUILD_CMD = CGO_ENABLED=0 go build -ldflags "-s -w"
BUILD_SRC = ./cmd/beatportdl
BUILD_DIR = ./bin

all: darwin-arm64 darwin-amd64 linux-amd64 linux-arm64 windows-amd64

darwin-arm64:
	@echo "Building for macOS ARM64"
	GOOS=darwin GOARCH=arm64 $(BUILD_CMD) -o=$(BUILD_DIR)/beatportdl-darwin-arm64 $(BUILD_SRC)

darwin-amd64:
	@echo "Building for macOS AMD64"
	GOOS=darwin GOARCH=amd64 $(BUILD_CMD) -o=$(BUILD_DIR)/beatportdl-darwin-amd64 $(BUILD_SRC)

linux-amd64:
	@echo "Building for Linux AMD64"
	GOOS=linux GOARCH=amd64 $(BUILD_CMD) -o=$(BUILD_DIR)/beatportdl-linux-amd64 $(BUILD_SRC)

linux-arm64:
	@echo "Building for Linux ARM64"
	GOOS=linux GOARCH=arm64 $(BUILD_CMD) -o=$(BUILD_DIR)/beatportdl-linux-arm64 $(BUILD_SRC)

windows-amd64:
	@echo "Building for Windows AMD64"
	GOOS=windows GOARCH=amd64 $(BUILD_CMD) -o=$(BUILD_DIR)/beatportdl-windows-amd64.exe $(BUILD_SRC)

.PHONY: all darwin-arm64 darwin-amd64 linux-amd64 linux-arm64 windows-amd64
