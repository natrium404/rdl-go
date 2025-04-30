#!/bin/bash

# For Linux build
wails build

# For Windows build
export GOOS=windows
export GOARCH=amd64
export WAILS_WINDOWS_USE_WINE=1
wails build
