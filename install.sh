#!/bin/bash
# SPIRIT Installation Script
# Usage: curl -fsSL https://theorionai.github.io/spirit/install.sh | bash

set -e

REPO="TheOrionAI/spirit"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case "$ARCH" in
        x86_64) ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
        *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
    esac
    
    case "$OS" in
        linux|darwin) ;;  # Supported
        *) echo "Unsupported OS: $OS"; exit 1 ;;
    esac
    
    PLATFORM="${OS}_${ARCH}"
}

# Get latest version
get_latest_version() {
    curl -s "https://api.github.com/repos/$REPO/releases/latest" | 
        grep '"tag_name":' | 
        sed -E 's/.*"([^"]+)".*/\1/' | 
        sed 's/^v//'
}

# Download and install
download() {
    local version="$1"
    local filename="spirit_${version}_${PLATFORM}.tar.gz"
    local url="https://github.com/$REPO/releases/download/v${version}/${filename}"
    
    echo "📥 Downloading SPIRIT v${version} for ${PLATFORM}..."
    echo "   URL: $url"
    
    # Create temp directory
    TMP_DIR=$(mktemp -d)
    trap "rm -rf $TMP_DIR" EXIT
    
    # Download
    if ! curl -fsSL "$url" -o "$TMP_DIR/spirit.tar.gz"; then
        echo "❌ Download failed"
        echo "   Make sure your platform is supported:"
        echo "   - Linux (amd64, arm64)"
        echo "   - macOS (amd64, arm64)"
        exit 1
    fi
    
    # Extract
    echo "📦 Extracting..."
    tar -xzf "$TMP_DIR/spirit.tar.gz" -C "$TMP_DIR"
    
    # Install
    echo "⚙️  Installing to $INSTALL_DIR..."
    if [ -w "$INSTALL_DIR" ]; then
        mv "$TMP_DIR/spirit" "$INSTALL_DIR/"
    else
        sudo mv "$TMP_DIR/spirit" "$INSTALL_DIR/"
    fi
    
    chmod +x "$INSTALL_DIR/spirit"
}

# Main
main() {
    echo "🌌 SPIRIT Installer"
    echo "===================="
    
    detect_platform
    VERSION=$(get_latest_version)
    
    download "$VERSION"
    
    # Verify
    if command -v spirit >/dev/null 2>&1; then
        echo ""
        echo "✅ SPIRIT installed successfully!"
        echo ""
        echo "   Version: $(spirit --version 2>/dev/null || echo 'v'$VERSION)"
        echo "   Location: $(which spirit)"
        echo ""
        echo "Quick start:"
        echo "  spirit init --name=\"my-agent\" --emoji=\"🤖\""
        echo "  spirit --help"
    else
        echo "⚠️  Installation complete, but 'spirit' not in PATH"
        echo "   Add $INSTALL_DIR to your PATH or run:"
        echo "   $INSTALL_DIR/spirit"
    fi
}

main "$@"
