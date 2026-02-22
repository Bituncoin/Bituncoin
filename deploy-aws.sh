#!/bin/bash

# BTNG Sovereign Gold Standard - AWS Deployment Script
# Builds Docker image, pushes to ECR, and deploys to ECS Fargate

set -e

# Configuration
AWS_REGION="us-east-1"
ECR_REPO_NAME="btng"
DOCKER_IMAGE_TAG="latest"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
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

check_dependencies() {
    log_info "Checking dependencies..."

    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed"
        exit 1
    fi

    if ! command -v aws &> /dev/null; then
        log_error "AWS CLI is not installed"
        exit 1
    fi

    if ! command -v terraform &> /dev/null; then
        log_error "Terraform is not installed"
        exit 1
    fi

    log_success "All dependencies found"
}

get_account_id() {
    aws sts get-caller-identity --query Account --output text
}

get_ecr_uri() {
    local account_id=$1
    echo "${account_id}.dkr.ecr.${AWS_REGION}.amazonaws.com/${ECR_REPO_NAME}"
}

build_docker_image() {
    log_info "Building Docker image..."

    # Build the multi-stage Docker image
    docker build -t btng:${DOCKER_IMAGE_TAG} .

    log_success "Docker image built successfully"
}

authenticate_ecr() {
    log_info "Authenticating with ECR..."

    local account_id=$(get_account_id)
    aws ecr get-login-password --region ${AWS_REGION} | \
        docker login --username AWS --password-stdin \
        ${account_id}.dkr.ecr.${AWS_REGION}.amazonaws.com

    log_success "Authenticated with ECR"
}

create_ecr_repository() {
    log_info "Ensuring ECR repository exists..."

    local account_id=$(get_account_id)

    # Check if repository exists
    if ! aws ecr describe-repositories --repository-names ${ECR_REPO_NAME} --region ${AWS_REGION} &> /dev/null; then
        log_info "Creating ECR repository..."
        aws ecr create-repository \
            --repository-name ${ECR_REPO_NAME} \
            --region ${AWS_REGION} \
            --image-scanning-configuration scanOnPush=true

        log_success "ECR repository created"
    else
        log_info "ECR repository already exists"
    fi
}

push_docker_image() {
    log_info "Pushing Docker image to ECR..."

    local ecr_uri=$(get_ecr_uri $(get_account_id))

    # Tag the image
    docker tag btng:${DOCKER_IMAGE_TAG} ${ecr_uri}:${DOCKER_IMAGE_TAG}

    # Push the image
    docker push ${ecr_uri}:${DOCKER_IMAGE_TAG}

    log_success "Docker image pushed to ECR: ${ecr_uri}:${DOCKER_IMAGE_TAG}"
}

setup_terraform() {
    log_info "Setting up Terraform..."

    # Initialize Terraform
    terraform init

    # Validate configuration
    terraform validate

    log_success "Terraform setup complete"
}

deploy_infrastructure() {
    log_info "Deploying infrastructure with Terraform..."

    # Plan the deployment
    terraform plan -var-file="terraform.tfvars" -out=tfplan

    # Apply the changes
    terraform apply tfplan

    log_success "Infrastructure deployed successfully"
}

run_integration_tests() {
    log_info "Running integration tests..."

    # Get the ALB DNS name
    local alb_dns=$(terraform output -raw alb_dns_name)

    if [ -z "$alb_dns" ]; then
        log_error "Could not get ALB DNS name from Terraform output"
        return 1
    fi

    log_info "Testing against ALB: ${alb_dns}"

    # Wait for services to be ready
    log_info "Waiting for services to be healthy..."
    sleep 120

    # Test health endpoint
    if curl -f -s "http://${alb_dns}/api/health" > /dev/null 2>&1; then
        log_success "Health check passed"
    else
        log_error "Health check failed"
        return 1
    fi

    # Test gold price status
    if curl -f -s "http://${alb_dns}/api/btng/gold/price/status" > /dev/null 2>&1; then
        log_success "Gold price status check passed"
    else
        log_error "Gold price status check failed"
        return 1
    fi

    log_success "All integration tests passed!"
}

cleanup() {
    log_info "Cleaning up temporary files..."

    # Remove Terraform plan
    rm -f tfplan

    # Remove local Docker images (optional)
    # docker rmi btng:${DOCKER_IMAGE_TAG} 2>/dev/null || true

    log_success "Cleanup complete"
}

main() {
    log_info "🚀 BTNG Sovereign Gold Standard - AWS Deployment"
    echo "=================================================="

    check_dependencies
    build_docker_image
    authenticate_ecr
    create_ecr_repository
    push_docker_image
    setup_terraform
    deploy_infrastructure
    run_integration_tests
    cleanup

    # Get outputs
    local alb_dns=$(terraform output -raw alb_dns_name)
    local ecr_uri=$(terraform output -raw ecr_repository_url)

    echo ""
    log_success "🎉 Deployment Complete!"
    echo ""
    echo "🌐 Application URL: http://${alb_dns}"
    echo "🔗 API Health Check: http://${alb_dns}/api/health"
    echo "📊 Gold Price Status: http://${alb_dns}/api/btng/gold/price/status"
    echo "🐳 ECR Repository: ${ecr_uri}"
    echo ""
    echo "📋 Next Steps:"
    echo "  1. Update your DNS to point to the ALB: ${alb_dns}"
    echo "  2. Configure HTTPS with AWS Certificate Manager"
    echo "  3. Set up monitoring and alerting"
    echo "  4. Configure auto-scaling based on your needs"
}

# Handle command line arguments
case "${1:-}" in
    "build")
        check_dependencies
        build_docker_image
        ;;
    "push")
        check_dependencies
        authenticate_ecr
        create_ecr_repository
        push_docker_image
        ;;
    "infra")
        setup_terraform
        deploy_infrastructure
        ;;
    "test")
        run_integration_tests
        ;;
    "cleanup")
        cleanup
        ;;
    *)
        main
        ;;
esac