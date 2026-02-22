#!/bin/bash

# BTNG Sovereign Platform - Production Deployment Script
# This script orchestrates the complete deployment process

set -e

echo "🚀 BTNG Sovereign Platform - Production Deployment"
echo "================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
NAMESPACE="btng-sovereign"
DOCKER_REGISTRY="${DOCKER_REGISTRY:-btng-sovereign}"
TAG="${TAG:-latest}"

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."

    # Check if kubectl is available
    if ! command -v kubectl &> /dev/null; then
        print_error "kubectl is not installed. Please install it first."
        exit 1
    fi

    # Check if docker is available
    if ! command -v docker &> /dev/null; then
        print_error "docker is not installed. Please install it first."
        exit 1
    fi

    # Check if kubernetes cluster is accessible
    if ! kubectl cluster-info &> /dev/null; then
        print_error "Cannot access Kubernetes cluster. Please check your kubeconfig."
        exit 1
    fi

    print_status "Prerequisites check passed."
}

# Build and push Docker images
build_and_push_images() {
    print_status "Building and pushing Docker images..."

    # Build BTNG SDK image
    print_status "Building BTNG SDK image..."
    docker build -f Dockerfile.btng-sdk -t ${DOCKER_REGISTRY}/btng-sdk:${TAG} .

    # Build Frontend image
    print_status "Building Frontend image..."
    docker build -f Dockerfile.frontend -t ${DOCKER_REGISTRY}/btng-frontend:${TAG} .

    # Push images
    print_status "Pushing images to registry..."
    docker push ${DOCKER_REGISTRY}/btng-sdk:${TAG}
    docker push ${DOCKER_REGISTRY}/btng-frontend:${TAG}

    print_status "Images built and pushed successfully."
}

# Deploy to Kubernetes
deploy_to_kubernetes() {
    print_status "Deploying to Kubernetes..."

    # Create namespace
    print_status "Creating namespace..."
    kubectl apply -f k8s/namespace.yml

    # Apply configurations
    print_status "Applying configurations..."
    kubectl apply -f k8s/configmap.yml
    kubectl apply -f k8s/secrets.yml
    kubectl apply -f k8s/pvc.yml

    # Deploy services
    print_status "Deploying services..."
    kubectl apply -f k8s/services.yml

    # Deploy applications
    print_status "Deploying applications..."
    kubectl apply -f k8s/sdk-deployment.yml
    kubectl apply -f k8s/frontend-deployment.yml

    # Apply networking and scaling
    print_status "Applying networking and scaling..."
    kubectl apply -f k8s/ingress.yml
    kubectl apply -f k8s/hpa.yml
    kubectl apply -f k8s/network-policy.yml

    print_status "Kubernetes deployment completed."
}

# Wait for rollout
wait_for_rollout() {
    print_status "Waiting for rollout to complete..."

    # Wait for SDK deployment
    kubectl rollout status deployment/btng-sdk -n ${NAMESPACE} --timeout=300s

    # Wait for Frontend deployment
    kubectl rollout status deployment/btng-frontend -n ${NAMESPACE} --timeout=300s

    print_status "Rollout completed successfully."
}

# Run health checks
run_health_checks() {
    print_status "Running health checks..."

    # Check pod status
    if kubectl get pods -n ${NAMESPACE} | grep -q "Running"; then
        print_status "All pods are running."
    else
        print_error "Some pods are not running. Check pod status."
        kubectl get pods -n ${NAMESPACE}
        exit 1
    fi

    # Check services
    if kubectl get svc -n ${NAMESPACE} | grep -q "btng-"; then
        print_status "Services are created."
    else
        print_error "Services are not created properly."
        exit 1
    fi

    print_status "Health checks passed."
}

# Main deployment function
main() {
    echo "Starting deployment process..."

    check_prerequisites
    build_and_push_images
    deploy_to_kubernetes
    wait_for_rollout
    run_health_checks

    echo ""
    print_status "🎉 Deployment completed successfully!"
    echo ""
    echo "Next steps:"
    echo "1. Update DNS records to point to your ingress IP"
    echo "2. Configure monitoring and alerting"
    echo "3. Set up backup procedures"
    echo "4. Configure SSL certificates"
    echo ""
    echo "Useful commands:"
    echo "  kubectl get pods -n ${NAMESPACE}"
    echo "  kubectl get svc -n ${NAMESPACE}"
    echo "  kubectl logs -n ${NAMESPACE} deployment/btng-sdk"
    echo "  kubectl logs -n ${NAMESPACE} deployment/btng-frontend"
}

# Run main function
main "$@"