#!/bin/sh
# Treehouse CLI installer
# Usage: curl -sSL https://get.treehouse.dev/install.sh | sh

set -e

# Check for required commands
command -v curl >/dev/null 2>&1 || { echo "✗ curl is required but not installed"; exit 1; }

REPO="jpmarques19/treehouse"
INSTALL_DIR="${HOME}/.local/bin"

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
  linux|darwin) ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
  x86_64|amd64)
    ARCH="amd64"
    ;;
  arm64|aarch64)
    ARCH="arm64"
    ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

BINARY="th-${OS}-${ARCH}"
echo "Detected: ${OS}/${ARCH}"

# Get latest release download URL
echo "Fetching latest release..."
LATEST_URL="https://api.github.com/repos/${REPO}/releases/latest"
DOWNLOAD_URL=$(curl -s "$LATEST_URL" | grep "browser_download_url.*${BINARY}\"" | head -1 | cut -d '"' -f 4)

if [ -z "$DOWNLOAD_URL" ]; then
  echo "Failed to find download URL for ${BINARY}"
  exit 1
fi

# Create install directory
mkdir -p "$INSTALL_DIR"

# Download and install
echo "Downloading ${BINARY}..."
curl -sL "$DOWNLOAD_URL" -o "${INSTALL_DIR}/th"
chmod +x "${INSTALL_DIR}/th"

echo ""
echo "✓ Installed th to ${INSTALL_DIR}/th"

# Display version
"${INSTALL_DIR}/th" --version 2>/dev/null || true

# Check if in PATH
case ":${PATH}:" in
  *":${INSTALL_DIR}:"*)
    echo ""
    echo "Run 'th --help' to get started"
    ;;
  *)
    echo ""
    echo "Add ${INSTALL_DIR} to your PATH:"
    echo "  export PATH=\"\${HOME}/.local/bin:\${PATH}\""
    ;;
esac
