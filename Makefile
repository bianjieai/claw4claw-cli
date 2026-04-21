.PHONY: all build build-linux-amd64 build-windows-amd64 build-linux-arm64 build-windows-arm64 build-linux-amd64-upx build-windows-amd64-upx build-linux-arm64-upx build-windows-arm64-upx clean

BINARY_NAME=c4c
LDFLAGS=-s -w
BUILD_FLAGS=-trimpath -ldflags="$(LDFLAGS)"

all: build

build:
	@echo "Building for macOS (Apple Silicon/ARM64)..."
	go build $(BUILD_FLAGS) -o $(BINARY_NAME) main.go

# --- Linux Builds ---
build-linux-amd64:
	@echo "Building for Linux (AMD64)..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-linux-amd64 main.go

build-linux-arm64:
	@echo "Building for Linux (ARM64)..."
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-linux-arm64 main.go

# --- Windows Builds ---
build-windows-amd64:
	@echo "Building for Windows (AMD64)..."
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-windows-amd64.exe main.go

build-windows-arm64:
	@echo "Building for Windows (ARM64)..."
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-windows-arm64.exe main.go

# --- UPX Compression ---
build-linux-amd64-upx: build-linux-amd64
	@echo "Compressing Linux AMD64 binary with UPX..."
	upx --best $(BINARY_NAME)-linux-amd64

build-linux-arm64-upx: build-linux-arm64
	@echo "Compressing Linux ARM64 binary with UPX..."
	upx --best $(BINARY_NAME)-linux-arm64

build-windows-amd64-upx: build-windows-amd64
	@echo "Compressing Windows AMD64 binary with UPX..."
	upx --best $(BINARY_NAME)-windows-amd64.exe

build-windows-arm64-upx: build-windows-arm64
	@echo "Compressing Windows ARM64 binary with UPX..."
	upx --best $(BINARY_NAME)-windows-arm64.exe

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME) $(BINARY_NAME)-*
