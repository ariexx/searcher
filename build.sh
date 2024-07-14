#!/bin/bash

# Build for Linux
echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o dist/searcher

# Build for Windows
echo "Building for Windows..."
GOOS=windows GOARCH=amd64 go build -o dist/searcher.exe

echo "Setting permissions..."
chmod +x dist/searcher

echo "Build complete."