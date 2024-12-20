#!/bin/bash

echo "Starting multi-platform build for wf..."

# Linux
echo "Building for Linux..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o build/wf-linux
if [ $? -eq 0 ]; then
    echo "Build successful for Linux: wf-linux"
else
    echo "Build failed for Linux"
    exit 1
fi

# Windows
echo "Building for Windows..."
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o build/wf-windows.exe
if [ $? -eq 0 ]; then
    echo "Build successful for Windows: wf-windows.exe"
else
    echo "Build failed for Windows"
    exit 1
fi

# MacOS
echo "Building for MacOS..."
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o build/wf-macos
if [ $? -eq 0 ]; then
    echo "Build successful for MacOS: wf-macos"
else
    echo "Build failed for MacOS"
    exit 1
fi

echo "All builds completed successfully! binaries are in the 'build' folder."