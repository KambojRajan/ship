#!/bin/bash

# Ship Remote Installation Script
# This script downloads and installs ship from GitHub releases or source

set -e

BINARY_NAME="ship"
INSTALL_PATH="/usr/local/bin"
REPO="KambojRajan/ship"
COLOR_GREEN='\033[0;32m'
COLOR_BLUE='\033[0;34m'
COLOR_RED='\033[0;31m'
COLOR_YELLOW='\033[0;33m'
COLOR_RESET='\033[0m'

echo -e "${COLOR_BLUE}🚢 Ship Remote Installation Script${COLOR_RESET}"
echo ""

# Detect OS and architecture
OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
    Linux*)     OS_TYPE="linux";;
    Darwin*)    OS_TYPE="darwin";;
    *)          OS_TYPE="unknown";;
esac

case "$ARCH" in
    x86_64)     ARCH_TYPE="amd64";;
    arm64)      ARCH_TYPE="arm64";;
    aarch64)    ARCH_TYPE="arm64";;
    *)          ARCH_TYPE="unknown";;
esac

echo -e "${COLOR_BLUE}📦 Detected platform: ${OS_TYPE}/${ARCH_TYPE}${COLOR_RESET}"

# Check if Go is installed (for source installation fallback)
if ! command -v go &> /dev/null; then
    echo -e "${COLOR_YELLOW}⚠️  Go is not installed. Will attempt to download pre-built binary.${COLOR_RESET}"
    HAS_GO=false
else
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    echo -e "${COLOR_BLUE}📦 Detected Go version: ${GO_VERSION}${COLOR_RESET}"
    HAS_GO=true
fi

# Try to get the latest release
echo -e "${COLOR_BLUE}🔍 Checking for latest release...${COLOR_RESET}"

# Attempt to download from GitHub releases
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' || echo "")

if [ -n "$LATEST_RELEASE" ]; then
    echo -e "${COLOR_BLUE}📥 Found release: ${LATEST_RELEASE}${COLOR_RESET}"
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_RELEASE}/ship-${OS_TYPE}-${ARCH_TYPE}"

    echo -e "${COLOR_BLUE}📥 Downloading binary...${COLOR_RESET}"
    if curl -L -f -o "/tmp/${BINARY_NAME}" "$DOWNLOAD_URL" 2>/dev/null; then
        echo -e "${COLOR_GREEN}✅ Download successful!${COLOR_RESET}"
        chmod +x "/tmp/${BINARY_NAME}"

        # Install the binary
        echo -e "${COLOR_BLUE}📦 Installing ${BINARY_NAME} to ${INSTALL_PATH}...${COLOR_RESET}"
        if [ -w "$INSTALL_PATH" ]; then
            mv "/tmp/${BINARY_NAME}" "$INSTALL_PATH/"
        else
            echo "Administrator privileges required for installation..."
            sudo mv "/tmp/${BINARY_NAME}" "$INSTALL_PATH/"
        fi

        echo ""
        echo -e "${COLOR_GREEN}✅ ${BINARY_NAME} installed successfully!${COLOR_RESET}"
        echo ""
        echo "Try running:"
        echo "  ${BINARY_NAME} --help"
        exit 0
    else
        echo -e "${COLOR_YELLOW}⚠️  Pre-built binary not available for ${OS_TYPE}/${ARCH_TYPE}${COLOR_RESET}"
    fi
fi

# Fallback: Install from source
if [ "$HAS_GO" = true ]; then
    echo -e "${COLOR_BLUE}🔨 Building from source...${COLOR_RESET}"

    # Create temporary directory
    TMP_DIR=$(mktemp -d)
    cd "$TMP_DIR"

    echo -e "${COLOR_BLUE}📥 Cloning repository...${COLOR_RESET}"
    git clone "https://github.com/${REPO}.git" . 2>/dev/null || {
        echo -e "${COLOR_RED}❌ Failed to clone repository${COLOR_RESET}"
        exit 1
    }

    echo -e "${COLOR_BLUE}📥 Downloading dependencies...${COLOR_RESET}"
    go mod download

    echo -e "${COLOR_BLUE}🔨 Building ${BINARY_NAME}...${COLOR_RESET}"
    go build -o "$BINARY_NAME" .

    if [ ! -f "$BINARY_NAME" ]; then
        echo -e "${COLOR_RED}❌ Build failed!${COLOR_RESET}"
        exit 1
    fi

    echo -e "${COLOR_GREEN}✅ Build successful!${COLOR_RESET}"

    # Install the binary
    echo -e "${COLOR_BLUE}📦 Installing ${BINARY_NAME} to ${INSTALL_PATH}...${COLOR_RESET}"
    if [ -w "$INSTALL_PATH" ]; then
        cp "$BINARY_NAME" "$INSTALL_PATH/"
        chmod +x "$INSTALL_PATH/$BINARY_NAME"
    else
        echo "Administrator privileges required for installation..."
        sudo cp "$BINARY_NAME" "$INSTALL_PATH/"
        sudo chmod +x "$INSTALL_PATH/$BINARY_NAME"
    fi

    # Cleanup
    cd - > /dev/null
    rm -rf "$TMP_DIR"

    echo ""
    echo -e "${COLOR_GREEN}✅ ${BINARY_NAME} installed successfully!${COLOR_RESET}"
    echo ""
    echo "Try running:"
    echo "  ${BINARY_NAME} --help"
    echo "  ${BINARY_NAME} init ."
else
    echo -e "${COLOR_RED}❌ Installation failed!${COLOR_RESET}"
    echo ""
    echo "Neither pre-built binaries nor Go compiler are available."
    echo "Please install Go and try again:"
    echo "  https://golang.org/dl/"
    exit 1
fi
