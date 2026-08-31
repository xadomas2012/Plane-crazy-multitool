#!/bin/sh
set -e

mkdir -p dist

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -trimpath \
  -ldflags="-s -w" \
  -o dist/PC-Gear-Calculator-Linux-x64

echo "Built: dist/PC-Gear-Calculator-Linux-x64"
