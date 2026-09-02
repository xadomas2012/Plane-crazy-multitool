#!/bin/sh
set -e

VERSION="${1:-dev}"

mkdir -p dist

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -trimpath \
  -ldflags="-s -w -X main.Version=${VERSION}" \
  -o dist/PC-Gear-Calculator-Linux-x64

echo "Built: dist/PC-Gear-Calculator-Linux-x64"
