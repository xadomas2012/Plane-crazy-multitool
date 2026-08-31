#!/bin/sh
set -e

mkdir -p dist

GOOS=windows GOARCH=amd64 go build \
    -trimpath \
    -ldflags="-s -w" \
    -o dist/PC-Gear-Calculator-Windows-x64.exe

echo "Built: dist/PC-Gear-Calculator-Windows-x64.exe"
