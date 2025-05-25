#!/bin/bash

set -e

# ==== Linux Build ====
echo "🔨 Building for Linux..."
export GOOS=linux
export GOARCH=amd64
wails build
LINUX_OUTPUT="rdl_linux_amd64"
mv build/bin/rdl build/bin/$LINUX_OUTPUT
echo "✅ Linux build: build/bin/$LINUX_OUTPUT"

# ==== Windows Build ====
echo "🔨 Building for Windows..."
export GOOS=windows
export GOARCH=amd64
export WAILS_WINDOWS_USE_WINE=1
wails build
WINDOWS_OUTPUT="rdl_windows_x64.exe"
mv build/bin/rdl.exe build/bin/$WINDOWS_OUTPUT
echo "✅ Windows build: build/bin/$WINDOWS_OUTPUT"
