#!/bin/bash
# RADAR installation script
# Created by Elite Security Systems (elitesecurity.systems)

set -e

# Detect OS and architecture
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

if [[ "$ARCH" == "x86_64" ]]; then
    ARCH="amd64"
elif [[ "$ARCH" == "aarch64" || "$ARCH" == "arm64" ]]; then
    ARCH="arm64"
else
    echo "Unsupported architecture: $ARCH"
    exit 1
fi

# Set version (can be overridden with RADAR_VERSION env var)
VERSION=${RADAR_VERSION:-latest}

if [[ "$VERSION" == "latest" ]]; then
    # Fetch the latest version from GitHub API
    if command -v curl > /dev/null; then
        VERSION=$(curl -s https://api.github.com/repos/Elite-Security-Systems/radar/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    elif command -v wget > /dev/null; then
        VERSION=$(wget -qO- https://api.github.com/repos/Elite-Security-Systems/radar/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    else
        echo "Neither curl nor wget found. Please install one of them or specify a version with RADAR_VERSION."
        exit 1
    fi
fi

# If we couldn't get the latest version, use a fallback
if [[ -z "$VERSION" ]]; then
    VERSION="v1.0.0"
    echo "Could not determine latest version, using $VERSION as fallback."
fi

# Set download URL
DOWNLOAD_URL="https://github.com/Elite-Security-Systems/radar/releases/download/$VERSION/radar-$OS-$ARCH"
if [[ "$OS" == "windows" ]]; then
    DOWNLOAD_URL="${DOWNLOAD_URL}.exe"
fi

# Set installation directory
if [[ "$OS" == "darwin" || "$OS" == "linux" ]]; then
    INSTALL_DIR="/usr/local/bin"
    if [[ ! -w "$INSTALL_DIR" ]]; then
        INSTALL_DIR="$HOME/.local/bin"
        mkdir -p "$INSTALL_DIR"
    fi
else
    INSTALL_DIR="$HOME/radar"
    mkdir -p "$INSTALL_DIR"
fi

# Set binary path
BINARY_PATH="$INSTALL_DIR/radar"
if [[ "$OS" == "windows" ]]; then
    BINARY_PATH="$INSTALL_DIR/radar.exe"
fi

echo "Installing RADAR $VERSION for $OS-$ARCH to $BINARY_PATH"

# Download the binary
if command -v curl > /dev/null; then
    curl -L "$DOWNLOAD_URL" -o "$BINARY_PATH"
elif command -v wget > /dev/null; then
    wget -O "$BINARY_PATH" "$DOWNLOAD_URL"
else
    echo "Neither curl nor wget found. Please install one of them."
    exit 1
fi

# Make the binary executable (on Unix-like systems)
if [[ "$OS" == "darwin" || "$OS" == "linux" ]]; then
    chmod +x "$BINARY_PATH"
fi

# Create signatures directory
if [[ "$OS" == "darwin" || "$OS" == "linux" ]]; then
    SIGNATURES_DIR="/usr/local/share/radar"
    if [[ ! -w /usr/local/share ]]; then
        SIGNATURES_DIR="$HOME/.local/share/radar"
    fi
else
    SIGNATURES_DIR="$INSTALL_DIR"
fi

mkdir -p "$SIGNATURES_DIR"

# Download the default signatures file
SIGNATURES_URL="https://raw.githubusercontent.com/Elite-Security-Systems/radar/$VERSION/data/signatures.json"
SIGNATURES_PATH="$SIGNATURES_DIR/signatures.json"

echo "Downloading default signatures to $SIGNATURES_PATH"

if command -v curl > /dev/null; then
    curl -L "$SIGNATURES_URL" -o "$SIGNATURES_PATH"
elif command -v wget > /dev/null; then
    wget -O "$SIGNATURES_PATH" "$SIGNATURES_URL"
fi

echo "Installation completed successfully!"
echo ""
echo "To use RADAR, run:"
echo "  radar -domain example.com"
echo ""
echo "For more options, run:"
echo "  radar -help"
