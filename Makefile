.PHONY: build install clean test help build-all release

BINARY_NAME=git-sync-remotes-tui
INSTALL_PATH=$(HOME)/.local/bin
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
DIST_DIR=dist

# Build flags
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

build:
	@echo "Building $(BINARY_NAME)..."
	@go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/tui
	@echo "Build complete: $(BINARY_NAME)"

install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	@mkdir -p $(INSTALL_PATH)
	@cp $(BINARY_NAME) $(INSTALL_PATH)/
	@chmod +x $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Installation complete!"
	@echo ""
	@echo "Make sure $(INSTALL_PATH) is in your PATH:"
	@echo "  export PATH=\"\$$PATH:\$$HOME/.local/bin\""

uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Uninstallation complete!"

clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)
	@rm -rf dist/
	@echo "Clean complete!"

test:
	@echo "Running tests..."
	@go test -v ./...

run: build
	@./$(BINARY_NAME)

fmt:
	@echo "Formatting code..."
	@go fmt ./...

lint:
	@echo "Running linter..."
	@go vet ./...

# Cross-platform builds
build-linux-amd64:
	@echo "Building for Linux (amd64)..."
	@mkdir -p $(DIST_DIR)
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/tui

build-linux-arm64:
	@echo "Building for Linux (arm64)..."
	@mkdir -p $(DIST_DIR)
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/tui

build-darwin-amd64:
	@echo "Building for macOS (Intel)..."
	@mkdir -p $(DIST_DIR)
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/tui

build-darwin-arm64:
	@echo "Building for macOS (Apple Silicon)..."
	@mkdir -p $(DIST_DIR)
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/tui

build-windows-amd64:
	@echo "Building for Windows (amd64)..."
	@mkdir -p $(DIST_DIR)
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/tui

build-all: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64
	@echo "All platform builds complete!"
	@ls -lh $(DIST_DIR)/

release: build-all
	@echo "Creating release archives..."
	@cd $(DIST_DIR) && \
		for binary in $(BINARY_NAME)-*; do \
			if [ -f "$$binary" ]; then \
				echo "Packaging $$binary..."; \
				tar czf "$$binary.tar.gz" "$$binary"; \
				shasum -a 256 "$$binary.tar.gz" > "$$binary.tar.gz.sha256"; \
			fi; \
		done
	@echo "Release artifacts created in $(DIST_DIR)/"
	@ls -lh $(DIST_DIR)/*.tar.gz

help:
	@echo "Available targets:"
	@echo "  build              - Build the binary for current platform"
	@echo "  install            - Build and install to ~/.local/bin"
	@echo "  uninstall          - Remove installed binary"
	@echo "  clean              - Remove build artifacts"
	@echo "  test               - Run tests"
	@echo "  run                - Build and run the TUI"
	@echo "  fmt                - Format code"
	@echo "  lint               - Run linter"
	@echo ""
	@echo "Cross-platform builds:"
	@echo "  build-linux-amd64  - Build for Linux (x86_64)"
	@echo "  build-linux-arm64  - Build for Linux (ARM64)"
	@echo "  build-darwin-amd64 - Build for macOS (Intel)"
	@echo "  build-darwin-arm64 - Build for macOS (Apple Silicon)"
	@echo "  build-windows-amd64- Build for Windows (x86_64)"
	@echo "  build-all          - Build for all platforms"
	@echo "  release            - Create release archives with checksums"
	@echo ""
	@echo "  help               - Show this help message"
