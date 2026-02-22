#!/bin/sh
# BTNG Sovereign Gold Standard - Docker Entrypoint Script
# Starts Nginx reverse proxy and both BTNG services

set -e

# Function to handle shutdown gracefully
shutdown() {
    echo "Shutting down services..."
    # Kill background processes
    kill $(jobs -p) 2>/dev/null || true
    exit 0
}

# Trap SIGTERM and SIGINT for graceful shutdown
trap shutdown SIGTERM SIGINT

echo "🚀 Starting BTNG Sovereign Gold Standard Platform"

# Start Genesis platform in background
echo "Starting Genesis platform on port 8080..."
java -jar /app/genesis.jar &
GENESIS_PID=$!

# Wait a moment for Genesis to start
sleep 10

# Start Next.js API in background
echo "Starting Next.js API on port 3003..."
cd /app/btng-api && node server.js &
API_PID=$!

# Wait a moment for API to start
sleep 5

# Start Nginx in foreground (this will be the main process)
echo "Starting Nginx reverse proxy on port 80..."
nginx -g "daemon off;" &
NGINX_PID=$!

# Wait for all background processes
wait $GENESIS_PID $API_PID $NGINX_PID