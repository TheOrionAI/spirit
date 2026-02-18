#!/bin/bash
set -e

# SPIRIT Installer
# 
# SECURITY NOTES:
# - This script downloads and installs the SPIRIT binary
# - Recommended: Download first, review, then execute:
#     curl -fsSL https://theorionai.github.io/spirit/install.sh -o /tmp/spirit-install.sh
#     cat /tmp/spirit-install.sh | head -50  # Review
#     bash /tmp/spirit-install.sh
#
# - Alternative: Use Homebrew
#     brew tap TheOrionAI/tap
#     brew install spirit
#
# - Alternative: Download directly from GitHub Releases

REPO="TheOrionAI/spirit"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Detect platform
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case "$ARCH" in
        x86_64) ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
        *) echo "❌ Unsupported architecture: $ARCH"; exit 1 ;;
    esac
    
    case "$OS" in
        linux|darwin) ;;
        *) echo "❌ Unsupported OS: $OS"; exit 1 ;;
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
        echo ""
        echo "❌ Download failed. Please check:"
        echo "   - Your internet connection"
        echo "   - Platform support (Linux/macOS, amd64/arm64)"
        echo ""
        echo "Alternatively, download manually from:"
        echo "   https://github.com/$REPO/releases"
        exit 1
    fi
    
    # Extract
    echo "📦 Extracting..."
    tar -xzf "$TMP_DIR/spirit.tar.gz" -C "$TMP_DIR"
    
    # Verify binary exists
    if [ ! -f "$TMP_DIR/spirit" ]; then
        echo "❌ Extracted archive doesn't contain 'spirit' binary"
        exit 1
    fi
    
    # Install
    echo "⚙️ Installing to $INSTALL_DIR..."
    if [ -w "$INSTALL_DIR" ]; then
        mv "$TMP_DIR/spirit" "$INSTALL_DIR/"
    else
        echo "   (requires sudo for $INSTALL_DIR)"
        sudo mv "$TMP_DIR/spirit" "$INSTALL_DIR/"
    fi
    
    chmod +x "$INSTALL_DIR/spirit"
}

# Main
main() {
    echo "🌌 SPIRIT Installer"
    echo "===================="
    echo ""
    
    # Platform detection
    detect_platform
    
    # Get version
    VERSION=$(get_latest_version)
    if [ -z "$VERSION" ]; then
        echo "❌ Failed to get latest version from GitHub API"
        echo "   You can download manually: https://github.com/$REPO/releases"
        exit 1
    fi
    
    # Download
    download "$VERSION"
    
    # Verify installation
    echo ""
    if command -v spirit >/dev/null 2>&1; then
        INSTALLED_VERSION=$(spirit --version 2>/dev/null || echo "v$VERSION")
        
        echo "✅ SPIRIT installed successfully!"
        echo ""
        echo "   Version: $INSTALLED_VERSION"
        echo "   Location: $(which spirit)"
        echo ""
        echo "Quick start:"
        echo "   spirit init --name=\"my-agent\" --emoji=\"🤖\""
        echo "   spirit --help"
        echo ""
        echo "Secure authentication (REQUIRED):"
        echo "   1. Create PRIVATE repo on GitHub"
        echo "   2. Use: gh auth login  OR  git config credential.helper store"
        echo "   3. NEVER: git remote add origin https://TOKEN@github.com/..."
        echo ""
    else
        echo "⚠️ Installation complete, but 'spirit' not in PATH"
        echo "   Add $INSTALL_DIR to PATH or run: $INSTALL_DIR/spirit"
    fi
}

main "$@"
