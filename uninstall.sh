#!/bin/bash

# git-sync-remotes uninstallation script
# This script removes git-sync-remotes installation

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
echo "git-sync-remotes Uninstallation"
echo "=========================================="
echo ""

# Default installation directory
INSTALL_DIR="${HOME}/.local/bin"
INSTALL_PATH="$INSTALL_DIR/git-sync-remotes"

# Check if installation exists
if [ ! -e "$INSTALL_PATH" ]; then
    print_warning "git-sync-remotes not found at $INSTALL_PATH"

    # Check for custom installation
    if command -v git-sync-remotes &> /dev/null; then
        CUSTOM_PATH=$(which git-sync-remotes)
        print_info "Found installation at: $CUSTOM_PATH"
        if confirm "Remove this installation?"; then
            INSTALL_PATH="$CUSTOM_PATH"
        else
            print_info "Uninstallation cancelled"
            exit 0
        fi
    else
        print_info "No installation found. Nothing to uninstall."
        exit 0
    fi
fi

# Show what will be removed
echo ""
print_info "The following will be removed:"
if [ -L "$INSTALL_PATH" ]; then
    echo "  - Symbolic link: $INSTALL_PATH"
    LINK_TARGET=$(readlink "$INSTALL_PATH")
    echo "    (points to: $LINK_TARGET)"
else
    echo "  - File: $INSTALL_PATH"
fi

echo ""
if ! confirm "Proceed with uninstallation?"; then
    print_info "Uninstallation cancelled"
    exit 0
fi

# Remove the installation
rm "$INSTALL_PATH"
print_success "Removed $INSTALL_PATH"

# Check if directory is now empty
if [ -d "$INSTALL_DIR" ] && [ -z "$(ls -A "$INSTALL_DIR")" ]; then
    echo ""
    print_info "$INSTALL_DIR is now empty"
    if confirm "Remove empty directory?"; then
        rmdir "$INSTALL_DIR"
        print_success "Removed directory $INSTALL_DIR"
    fi
fi

echo ""
echo "=========================================="
print_success "Uninstallation complete!"
echo "=========================================="
echo ""
print_info "git-sync-remotes has been removed from your system"

# Check if the source repository still exists
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
if [ -d "$SCRIPT_DIR" ]; then
    echo ""
    print_info "The source repository is still available at:"
    echo "  $SCRIPT_DIR"
    echo ""
    print_info "To reinstall, run:"
    echo "  cd $SCRIPT_DIR && ./install.sh"
fi
