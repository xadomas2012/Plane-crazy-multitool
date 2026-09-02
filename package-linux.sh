#!/bin/sh
set -eu

VERSION="${1:-dev}"

PACKAGE_DIR="dist/linux-package"
APP_DIR="$PACKAGE_DIR/PC-Multitool"
ZIP="dist/PC-Multitool-Linux-x64-v${VERSION}.zip"

rm -rf "$APP_DIR"
mkdir -p "$APP_DIR"

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -ldflags="-s -w -X main.Version=${VERSION}" \
    -o "$APP_DIR/PC-Gear-Calculator"

cp install.sh "$APP_DIR/install.sh"

chmod +x \
    "$APP_DIR/PC-Gear-Calculator" \
    "$APP_DIR/install.sh"

rm -f "$ZIP"

(
    cd "$PACKAGE_DIR"
    zip -qr "../../$ZIP" PC-Multitool
)

echo "Built: $ZIP"
