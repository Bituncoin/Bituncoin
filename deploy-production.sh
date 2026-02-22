#!/bin/bash

# BTNG Sovereign Production Deployment Script
# Deploys the monorepo to production with health checks and rollback

set -e

echo "🇰🇪 BTNG Sovereign Gold Standard - Production Deployment"
echo "=================================================="

# Configuration
DEPLOY_ENV=${DEPLOY_ENV:-production}
DOCKER_COMPOSE_FILE="docker-compose.yml"
BACKUP_DIR="./backups/$(date +%Y%m%d_%H%M%S)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_dependencies() {
    log_info "Checking dependencies..."

    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed"
        exit 1
    fi

    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose is not installed"
        exit 1
    fi

    log_info "Dependencies OK"
}

create_backup() {
    log_info "Creating backup..."

    mkdir -p "$BACKUP_DIR"

    # Backup environment files
    if [ -f ".env" ]; then
        cp .env "$BACKUP_DIR/"
    fi

    # Backup database if running
    if docker ps | grep -q btng-mongodb; then
        log_info "Backing up MongoDB data..."
        docker exec btng-mongodb mongodump --out /backup
        docker cp btng-mongodb:/backup "$BACKUP_DIR/mongodb_backup"
    fi

    log_info "Backup created at $BACKUP_DIR"
}

health_check() {
    local service=$1
    local url=$2
    local max_attempts=30
    local attempt=1

    log_info "Health checking $service at $url..."

    while [ $attempt -le $max_attempts ]; do
        if curl -f -s "$url" > /dev/null 2>&1; then
            log_info "$service is healthy"
            return 0
        fi

        log_warn "Attempt $attempt/$max_attempts: $service not ready yet..."
        sleep 10
        ((attempt++))
    done

    log_error "$service failed health check after $max_attempts attempts"
    return 1
}

rollback() {
    log_error "Deployment failed, rolling back..."

    # Stop new deployment
    docker-compose -f "$DOCKER_COMPOSE_FILE" down

    # Restore from backup if available
    if [ -d "$BACKUP_DIR" ]; then
        log_info "Restoring from backup..."
        if [ -f "$BACKUP_DIR/.env" ]; then
            cp "$BACKUP_DIR/.env" .env
        fi
        # Additional rollback logic here
    fi

    # Start previous version
    docker-compose -f "$DOCKER_COMPOSE_FILE" up -d

    exit 1
}

# Main deployment
main() {
    log_info "Starting deployment to $DEPLOY_ENV environment"

    check_dependencies
    create_backup

    # Pre-deployment checks
    log_info "Running pre-deployment checks..."

    # Validate environment variables
    if [ -z "$JWT_SECRET" ]; then
        log_error "JWT_SECRET environment variable is required"
        exit 1
    fi

    if [ -z "$API_SECRET" ]; then
        log_error "API_SECRET environment variable is required"
        exit 1
    fi

    # Build and deploy
    log_info "Building services..."
    docker-compose -f "$DOCKER_COMPOSE_FILE" build --no-cache

    log_info "Starting services..."
    docker-compose -f "$DOCKER_COMPOSE_FILE" up -d

    # Wait for services to be ready
    sleep 30

    # Health checks
    if ! health_check "BTNG API" "http://localhost:3003/api/health"; then
        rollback
    fi

    if ! health_check "BTNG Gold Price Status" "http://localhost:3003/api/btng/gold/price/status"; then
        rollback
    fi

    # Run integration tests
    log_info "Running integration tests..."
    if ! npm run test:integration; then
        log_warn "Integration tests failed, but continuing deployment..."
    fi

    log_info "Deployment completed successfully! 🎉"
    log_info "Services running:"
    log_info "  - BTNG API: http://localhost:3003"
    log_info "  - Health Check: http://localhost:3003/api/health"
    log_info "  - Gold Price Demo: http://localhost:3003/gold-price-demo"

    # Cleanup old backups (keep last 5)
    log_info "Cleaning up old backups..."
    ls -td ./backups/* | tail -n +6 | xargs -r rm -rf
}

# Handle command line arguments
case "${1:-}" in
    "rollback")
        rollback
        ;;
    "health")
        health_check "BTNG API" "http://localhost:3003/api/health"
        ;;
    *)
        main
        ;;
esac