#!/bin/bash
# Test script for BTNG Docker container with Nginx

echo "🧪 Testing BTNG Docker container with Nginx reverse proxy"

# Run the container in background
echo "Starting container..."
docker run -d --name btng-test -p 8080:80 btng:latest

# Wait for services to start
echo "Waiting for services to start..."
sleep 30

# Test health endpoint
echo "Testing health endpoint..."
curl -f http://localhost:8080/health
if [ $? -eq 0 ]; then
    echo "✅ Health check passed"
else
    echo "❌ Health check failed"
fi

# Test API proxy (this will fail since we don't have the actual API running, but nginx should respond)
echo "Testing API proxy..."
curl -s http://localhost:8080/api/health || echo "Expected - API not running in test"

# Test Genesis proxy (this will fail since we don't have Genesis running, but nginx should respond)
echo "Testing Genesis proxy..."
curl -s http://localhost:8080/genesis/health || echo "Expected - Genesis not running in test"

# Check nginx configuration
echo "Checking nginx configuration..."
docker exec btng-test nginx -t
if [ $? -eq 0 ]; then
    echo "✅ Nginx configuration is valid"
else
    echo "❌ Nginx configuration is invalid"
fi

# Check running processes
echo "Checking running processes..."
docker exec btng-test ps aux

# Clean up
echo "Cleaning up..."
docker stop btng-test
docker rm btng-test

echo "🎉 Test completed!"