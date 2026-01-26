#!/bin/bash

# Ship Installation Script
# This script automatically builds and installs the ship binary

set -e

BINARY_NAME="ship"
INSTALL_PATH="/usr/local/bin"
BUILD_DIR="bin"
COLOR_GREEN='\033[0;32m'
COLOR_BLUE='\033[0;34m'
COLOR_RED='\033[0;31m'
COLOR_RESET='\033[0m'

echo -e "${COLOR_BLUE}🚢 Ship Installation Script${COLOR_RESET}"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${COLOR_RED}❌ Go is not installed. Please install Go 1.24.10 or higher.${COLOR_RESET}"
    echo "Visit: https://golang.org/dl/"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo -e "${COLOR_BLUE}📦 Detected Go version: ${GO_VERSION}${COLOR_RESET}"

# Download dependencies
echo -e "${COLOR_BLUE}📥 Downloading dependencies...${COLOR_RESET}"
go mod download

# Build the binary
echo -e "${COLOR_BLUE}🔨 Building ${BINARY_NAME}...${COLOR_RESET}"
mkdir -p "$BUILD_DIR"
go build -o "$BUILD_DIR/$BINARY_NAME" .

# Check if build was successful
if [ ! -f "$BUILD_DIR/$BINARY_NAME" ]; then
    echo -e "${COLOR_RED}❌ Build failed!${COLOR_RESET}"
    exit 1
fi

echo -e "${COLOR_GREEN}✅ Build successful!${COLOR_RESET}"

# Install the binary
echo -e "${COLOR_BLUE}📦 Installing ${BINARY_NAME} to ${INSTALL_PATH}...${COLOR_RESET}"

if [ -w "$INSTALL_PATH" ]; then
    # User has write permission
    cp "$BUILD_DIR/$BINARY_NAME" "$INSTALL_PATH/"
    chmod +x "$INSTALL_PATH/$BINARY_NAME"
else
    # Need sudo
    echo "Administrator privileges required for installation..."
    sudo cp "$BUILD_DIR/$BINARY_NAME" "$INSTALL_PATH/"
    sudo chmod +x "$INSTALL_PATH/$BINARY_NAME"
fi

# Verify installation
if command -v "$BINARY_NAME" &> /dev/null; then
    echo ""
    echo -e "${COLOR_GREEN}✅ ${BINARY_NAME} installed successfully!${COLOR_RESET}"
    echo ""
    echo -e "${COLOR_BLUE}You can now use '${BINARY_NAME}' from anywhere.${COLOR_RESET}"
    echo ""
    echo "Try running:"
    echo "  ${BINARY_NAME} --help"
    echo "  ${BINARY_NAME} init ."
else
    echo -e "${COLOR_RED}❌ Installation verification failed.${COLOR_RESET}"
    echo "The binary was copied but is not in PATH."
    echo "You may need to add ${INSTALL_PATH} to your PATH variable."
    exit 1
fi
