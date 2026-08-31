#!/bin/sh
set -e

mkdir -p dist

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build \
    -trimpath \
    -ldflags="-s -w" \
    -o dist/PC-Gear-Calculator-macOS-x64

CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build \
    -trimpath \
    -ldflags="-s -w" \
    -o dist/PC-Gear-Calculator-macOS-arm64

echo "Built:"
echo "  dist/PC-Gear-Calculator-macOS-x64"
echo "  dist/PC-Gear-Calculator-macOS-arm64"
