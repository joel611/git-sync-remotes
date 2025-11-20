.PHONY: build install clean test help

BINARY_NAME=git-sync-remotes-tui
INSTALL_PATH=$(HOME)/.local/bin

build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) ./cmd/tui
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

help:
	@echo "Available targets:"
	@echo "  build      - Build the binary"
	@echo "  install    - Build and install to ~/.local/bin"
	@echo "  uninstall  - Remove installed binary"
	@echo "  clean      - Remove build artifacts"
	@echo "  test       - Run tests"
	@echo "  run        - Build and run the TUI"
	@echo "  fmt        - Format code"
	@echo "  lint       - Run linter"
	@echo "  help       - Show this help message"
