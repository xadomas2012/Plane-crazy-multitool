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

mkdir -p "$APP_DIR/.update"

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -ldflags="-s -w" \
    -o "$APP_DIR/.update/PC-Gear-Calculator-Updater" \
    ./cmd/updater

chmod +x \
    "$APP_DIR/PC-Gear-Calculator" \
    "$APP_DIR/install.sh" \
    "$APP_DIR/.update/PC-Gear-Calculator-Updater"

rm -f "$ZIP"

(
    cd "$PACKAGE_DIR"
    zip -qr "../../$ZIP" PC-Multitool
)

echo "Built: $ZIP"
