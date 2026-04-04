#!/bin/sh
set -e

# Start .NET BTNG Wallet API on port 64799
cd /app/dotnet
dotnet McpServer.dll --urls "http://0.0.0.0:64799" &

# Start Next.js frontend on port 3000
cd /app
npx next start btng-api -p 3000 -H 0.0.0.0
