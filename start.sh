#!/bin/bash

# Bituncoin Universal Wallet - Quick Start Script

echo "🚀 Bituncoin Universal Wallet - Quick Start"
echo "==========================================="
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

echo "✓ Go $(go version | awk '{print $3}') detected"

# Navigate to bituncoin-btn directory
cd bituncoin-btn

# Download Go dependencies
echo ""
echo "📦 Downloading Go dependencies..."
go mod download

# Build the node
echo ""
echo "🔨 Building Bituncoin node..."
go build -o bituncoin-node main.go

if [ $? -eq 0 ]; then
    echo "✓ Build successful!"
    echo ""
    echo "🎉 Setup complete! You can now run the node with:"
    echo "   cd bituncoin-btn && ./bituncoin-node"
    echo ""
    echo "📚 For wallet UI setup, see the README.md"
else
    echo "❌ Build failed. Please check the errors above."
    exit 1
fi
