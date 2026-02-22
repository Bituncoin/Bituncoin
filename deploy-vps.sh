# BTNG VPS Deployment Script
# Deploys the BTNG sovereign gold standard platform to Ubuntu VPS with Docker

#!/bin/bash

set -e

echo "🚀 BTNG Sovereign Gold Standard - VPS Deployment"
echo "=================================================="

# Configuration
DOCKER_IMAGE="btng:latest"
CONTAINER_NAME="btng-app"
HOST_PORT=80

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running on Ubuntu with Docker
check_system() {
    log_info "Checking system requirements..."

    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed"
        exit 1
    fi

    if ! command -v node &> /dev/null; then
        log_error "Node.js is not installed"
        exit 1
    fi

    if ! command -v npm &> /dev/null; then
        log_error "npm is not installed"
        exit 1
    fi

    log_success "System requirements met"
}

# Build the Next.js application locally
build_nextjs() {
    log_info "Building Next.js application..."

    cd btng-api

    # Clean previous build
    rm -rf .next

    # Install dependencies
    npm install

    # Build the application
    npm run build

    if [ $? -ne 0 ]; then
        log_error "Next.js build failed"
        exit 1
    fi

    cd ..
    log_success "Next.js application built successfully"
}

# Build the Genesis JAR
build_genesis() {
    log_info "Building Genesis platform..."

    cd genesis-app

    # Build the JAR
    ./gradlew clean build --no-daemon -x test

    if [ $? -ne 0 ]; then
        log_error "Genesis build failed"
        exit 1
    fi

    cd ..
    log_success "Genesis platform built successfully"
}

# Build Docker image
build_docker() {
    log_info "Building Docker image..."

    docker build -t $DOCKER_IMAGE .

    if [ $? -ne 0 ]; then
        log_error "Docker build failed"
        exit 1
    fi

    log_success "Docker image built successfully"
}

# Deploy the container
deploy_container() {
    log_info "Deploying BTNG container..."

    # Stop and remove existing container if it exists
    docker stop $CONTAINER_NAME 2>/dev/null || true
    docker rm $CONTAINER_NAME 2>/dev/null || true

    # Run the new container
    docker run -d \
        --name $CONTAINER_NAME \
        --restart unless-stopped \
        -p $HOST_PORT:80 \
        --env-file .env \
        $DOCKER_IMAGE

    if [ $? -ne 0 ]; then
        log_error "Container deployment failed"
        exit 1
    fi

    log_success "BTNG container deployed successfully"
}

# Health check
health_check() {
    log_info "Performing health checks..."

    # Wait for container to start
    sleep 10

    # Check if container is running
    if ! docker ps | grep -q $CONTAINER_NAME; then
        log_error "Container is not running"
        exit 1
    fi

    # Check health endpoint
    max_attempts=10
    attempt=1

    while [ $attempt -le $max_attempts ]; do
        log_info "Health check attempt $attempt/$max_attempts..."

        if curl -f -s http://localhost:$HOST_PORT/health > /dev/null 2>&1; then
            log_success "Health check passed!"
            return 0
        fi

        sleep 5
        ((attempt++))
    done

    log_error "Health check failed after $max_attempts attempts"
    log_info "Container logs:"
    docker logs $CONTAINER_NAME
    exit 1
}

# Show deployment info
show_info() {
    log_success "🎉 BTNG Deployment Complete!"
    echo ""
    echo "🌐 Application URL: http://$(hostname -I | awk '{print $1}'):$HOST_PORT"
    echo "🔗 Health Check: http://$(hostname -I | awk '{print $1}'):$HOST_PORT/health"
    echo "🔗 API Health Check: http://$(hostname -I | awk '{print $1}'):$HOST_PORT/api/health"
    echo "🐳 Container Name: $CONTAINER_NAME"
    echo ""
    echo "📋 Useful commands:"
    echo "  docker logs $CONTAINER_NAME          # View logs"
    echo "  docker restart $CONTAINER_NAME       # Restart container"
    echo "  docker exec -it $CONTAINER_NAME sh   # Access container shell"
    echo ""
}

# Main deployment flow
main() {
    check_system
    build_nextjs
    build_genesis
    build_docker
    deploy_container
    health_check
    show_info
}

# Handle command line arguments
case "${1:-all}" in
    "build")
        check_system
        build_nextjs
        build_genesis
        build_docker
        ;;
    "deploy")
        check_system
        deploy_container
        ;;
    "health")
        health_check
        ;;
    "info")
        show_info
        ;;
    "all")
        main
        ;;
    *)
        echo "Usage: $0 [build|deploy|health|info|all]"
        echo "  build  - Build applications and Docker image"
        echo "  deploy - Deploy the container"
        echo "  health - Run health checks"
        echo "  info   - Show deployment information"
        echo "  all    - Full deployment (default)"
        exit 1
        ;;
esac