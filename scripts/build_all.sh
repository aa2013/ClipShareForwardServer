#!/bin/bash
set -e

VERSION=$1

if [ -z "$VERSION" ]; then
  echo "Version not provided"
  exit 1
fi

echo "Building version: $VERSION"

# =========================
# Build front
# =========================
npm --prefix ./front install
npm --prefix ./front run build

rm -rf ./server/web/dist
cp -r ./front/dist ./server/web

# =========================
# Build server
# =========================
mkdir -p build/linux
mkdir -p build/windows

# =========================
# Linux amd64
# =========================
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 \
go build -C ./server -ldflags="-s -w" \
-o ../build/linux/forward_server_linux_amd64_${VERSION} ./

# =========================
# Linux arm64
# =========================
GOOS=linux GOARCH=arm64 CGO_ENABLED=1 \
CC=aarch64-linux-gnu-gcc \
CXX=aarch64-linux-gnu-g++ \
go build -C ./server -ldflags="-s -w" \
-o ../build/linux/forward_server_linux_arm64_${VERSION} ./

# =========================
# Windows amd64
# =========================
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 \
CC=x86_64-w64-mingw32-gcc \
CXX=x86_64-w64-mingw32-g++ \
go build -C ./server -ldflags="-s -w" \
-o ../build/windows/forward_server_windows_amd64_${VERSION}.exe ./

echo "Build complete"