#!/bin/bash
# RADAR build script
# Created by Elite Security Systems (elitesecurity.systems)

set -e

# Get version information
VERSION=$(git describe --tags --always 2>/dev/null || echo "dev")
COMMIT=$(git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ')
LDFLAGS="-X main.Version=$VERSION -X main.Commit=$COMMIT -X main.BuildDate=$BUILD_DATE"

# Build for the current platform
echo "Building RADAR version $VERSION (commit $COMMIT)"

# Parse command line arguments
BUILD_ALL=false
OUTPUT_DIR="./bin"
PLATFORMS=()

while [[ $# -gt 0 ]]; do
  case $1 in
    --all)
      BUILD_ALL=true
      shift
      ;;
    --output)
      OUTPUT_DIR="$2"
      shift 2
      ;;
    --platform)
      PLATFORMS+=("$2")
      shift 2
      ;;
    *)
      echo "Unknown option: $1"
      echo "Usage: $0 [--all] [--output <dir>] [--platform <os>-<arch>]"
      exit 1
      ;;
  esac
done

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Define build platforms
if [[ "$BUILD_ALL" == "true" ]]; then
  PLATFORMS=(
    "linux-amd64"
    "linux-arm64"
    "darwin-amd64"
    "darwin-arm64"
    "windows-amd64"
  )
elif [[ ${#PLATFORMS[@]} -eq 0 ]]; then
  # Detect current platform
  OS=$(uname -s | tr '[:upper:]' '[:lower:]')
  ARCH=$(uname -m)
  
  if [[ "$ARCH" == "x86_64" ]]; then
    ARCH="amd64"
  elif [[ "$ARCH" == "aarch64" || "$ARCH" == "arm64" ]]; then
    ARCH="arm64"
  fi
  
  PLATFORMS=("$OS-$ARCH")
fi

# Build for each platform
for platform in "${PLATFORMS[@]}"; do
  os="${platform%%-*}"
  arch="${platform##*-}"
  output="$OUTPUT_DIR/radar-$os-$arch"
  
  if [[ "$os" == "windows" ]]; then
    output="$output.exe"
  fi
  
  echo "Building for $os/$arch -> $output"
  GOOS=$os GOARCH=$arch go build -ldflags "$LDFLAGS" -o "$output" ./cmd/radar
  
  # Create checksums
  if command -v sha256sum > /dev/null; then
    sha256sum "$output" > "$output.sha256"
  elif command -v shasum > /dev/null; then
    shasum -a 256 "$output" > "$output.sha256"
  fi
done

echo "Build completed successfully!"
echo "Binaries are available in $OUTPUT_DIR"
