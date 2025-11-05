#!/bin/bash

# git-sync-remotes installation script
# This script installs git-sync-remotes by creating a symbolic link

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_error() {
    echo -e "${RED}ERROR: $1${NC}" >&2
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

confirm() {
    local prompt="$1"
    local response
    read -p "$prompt [y/n]: " response
    case "$response" in
        [Yy]* ) return 0;;
        [Nn]* ) return 1;;
        * )
            print_warning "Invalid response, treating as 'no'"
            return 1
            ;;
    esac
}

echo "=========================================="
echo "git-sync-remotes Installation"
echo "=========================================="
echo ""

# Get the directory where this script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
SOURCE_SCRIPT="$SCRIPT_DIR/git-sync-remotes"

# Check if source script exists
if [ ! -f "$SOURCE_SCRIPT" ]; then
    print_error "git-sync-remotes script not found in $SCRIPT_DIR"
    exit 1
fi

# Make source script executable
chmod +x "$SOURCE_SCRIPT"
print_success "Made script executable"

# Default installation directory
INSTALL_DIR="${HOME}/.local/bin"

# Ask user for installation directory
print_info "Default installation directory: $INSTALL_DIR"
if confirm "Use default directory?"; then
    mkdir -p "$INSTALL_DIR"
else
    read -p "Enter installation directory: " custom_dir
    INSTALL_DIR="${custom_dir/#\~/$HOME}"
    mkdir -p "$INSTALL_DIR"
fi

INSTALL_PATH="$INSTALL_DIR/git-sync-remotes"

# Remove existing file/link if it exists
if [ -e "$INSTALL_PATH" ]; then
    if [ -L "$INSTALL_PATH" ]; then
        print_warning "Existing symbolic link found at $INSTALL_PATH"
    else
        print_warning "Existing file found at $INSTALL_PATH"
    fi

    if confirm "Remove existing installation?"; then
        rm "$INSTALL_PATH"
        print_success "Removed existing installation"
    else
        print_error "Installation cancelled"
        exit 1
    fi
fi

# Create symbolic link
ln -s "$SOURCE_SCRIPT" "$INSTALL_PATH"
print_success "Created symbolic link: $INSTALL_PATH -> $SOURCE_SCRIPT"

# Check if directory is in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    print_warning "$INSTALL_DIR is not in your PATH"
    echo ""
    print_info "Add the following to your shell config file (~/.bashrc or ~/.zshrc):"
    echo "  export PATH=\"\$PATH:$INSTALL_DIR\""
    echo ""
    print_info "Then reload your shell:"
    echo "  source ~/.zshrc  # or ~/.bashrc"
fi

echo ""
echo "=========================================="
print_success "Installation complete!"
echo "=========================================="
echo ""
print_info "Usage examples:"
echo "  git-sync-remotes        # Sync current branch"
echo "  git-sync-remotes -y     # Sync with auto-confirm"
echo "  git-sync-remotes master # Sync master branch"
echo ""
print_info "For more information, see: $SCRIPT_DIR/README.md"
