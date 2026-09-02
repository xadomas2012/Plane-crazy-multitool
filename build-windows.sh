#!/bin/sh
set -e

VERSION="${1:-dev}"

mkdir -p dist

GOOS=windows GOARCH=amd64 go build \
    -trimpath \
    -ldflags="-s -w -X main.Version=${VERSION}" \
    -o dist/PC-Gear-Calculator-Windows-x64.exe

echo "Built: dist/PC-Gear-Calculator-Windows-x64.exe"
