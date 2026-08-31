#!/bin/sh
set -e

mkdir -p dist

go build \
    -trimpath \
    -ldflags="-s -w" \
    -o dist/PC-Gear-Calculator-Linux-x64

echo "Built: dist/PC-Gear-Calculator-Linux-x64"
