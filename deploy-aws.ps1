# BTNG Sovereign Gold Standard - AWS Deployment Script
# Builds Docker image, pushes to ECR, and deploys to ECS Fargate

param(
    [string]$Action = "all"
)

# Configuration
$AWS_REGION = "us-east-1"
$ECR_REPO_NAME = "btng"
$DOCKER_IMAGE_TAG = "latest"

# Colors for output
$RED = [char]27 + "[31m"
$GREEN = [char]27 + "[32m"
$YELLOW = [char]27 + "[33m"
$BLUE = [char]27 + "[34m"
$NC = [char]27 + "[0m" # No Color

# Functions
function Write-Info {
    param([string]$Message)
    Write-Host "$BLUE[INFO]$NC $Message"
}

function Write-Success {
    param([string]$Message)
    Write-Host "$GREEN[SUCCESS]$NC $Message"
}

function Write-Warn {
    param([string]$Message)
    Write-Host "$YELLOW[WARN]$NC $Message"
}

function Write-Error {
    param([string]$Message)
    Write-Host "$RED[ERROR]$NC $Message"
}

function Test-Dependencies {
    Write-Info "Checking dependencies..."

    if (!(Get-Command docker -ErrorAction SilentlyContinue)) {
        Write-Error "Docker is not installed"
        exit 1
    }

    if (!(Get-Command aws -ErrorAction SilentlyContinue)) {
        Write-Error "AWS CLI is not installed"
        exit 1
    }

    if (!(Get-Command terraform -ErrorAction SilentlyContinue)) {
        Write-Error "Terraform is not installed"
        exit 1
    }

    Write-Success "All dependencies found"
}

function Get-AccountId {
    return aws sts get-caller-identity --query Account --output text
}

function Get-EcrUri {
    param([string]$AccountId)
    return "$AccountId.dkr.ecr.$AWS_REGION.amazonaws.com/$ECR_REPO_NAME"
}

function Build-DockerImage {
    Write-Info "Building Docker image..."

    # Build the multi-stage Docker image
    docker build -t btng:$DOCKER_IMAGE_TAG .

    if ($LASTEXITCODE -ne 0) {
        Write-Error "Docker build failed"
        exit 1
    }

    Write-Success "Docker image built successfully"
}

function Authenticate-Ecr {
    Write-Info "Authenticating with ECR..."

    $accountId = Get-AccountId
    aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin "$accountId.dkr.ecr.$AWS_REGION.amazonaws.com"

    if ($LASTEXITCODE -ne 0) {
        Write-Error "ECR authentication failed"
        exit 1
    }

    Write-Success "Authenticated with ECR"
}

function New-EcrRepository {
    Write-Info "Ensuring ECR repository exists..."

    $accountId = Get-AccountId

    # Check if repository exists
    $repoExists = aws ecr describe-repositories --repository-names $ECR_REPO_NAME --region $AWS_REGION 2>$null

    if ($LASTEXITCODE -ne 0) {
        Write-Info "Creating ECR repository..."
        aws ecr create-repository `
            --repository-name $ECR_REPO_NAME `
            --region $AWS_REGION `
            --image-scanning-configuration scanOnPush=true

        if ($LASTEXITCODE -ne 0) {
            Write-Error "Failed to create ECR repository"
            exit 1
        }

        Write-Success "ECR repository created"
    } else {
        Write-Info "ECR repository already exists"
    }
}

function Push-DockerImage {
    Write-Info "Pushing Docker image to ECR..."

    $accountId = Get-AccountId
    $ecrUri = Get-EcrUri $accountId

    # Tag the image
    docker tag btng:$DOCKER_IMAGE_TAG "$ecrUri`:$DOCKER_IMAGE_TAG"

    if ($LASTEXITCODE -ne 0) {
        Write-Error "Failed to tag Docker image"
        exit 1
    }

    # Push the image
    docker push "$ecrUri`:$DOCKER_IMAGE_TAG"

    if ($LASTEXITCODE -ne 0) {
        Write-Error "Failed to push Docker image"
        exit 1
    }

    Write-Success "Docker image pushed to ECR: $ecrUri`:$DOCKER_IMAGE_TAG"
}

function Setup-Terraform {
    Write-Info "Setting up Terraform..."

    # Initialize Terraform
    terraform init

    if ($LASTEXITCODE -ne 0) {
        Write-Error "Terraform init failed"
        exit 1
    }

    # Validate configuration
    terraform validate

    if ($LASTEXITCODE -ne 0) {
        Write-Error "Terraform validation failed"
        exit 1
    }

    Write-Success "Terraform setup complete"
}

function Deploy-Infrastructure {
    Write-Info "Deploying infrastructure with Terraform..."

    # Plan the deployment
    terraform plan -var-file="terraform.tfvars" -out=tfplan

    if ($LASTEXITCODE -ne 0) {
        Write-Error "Terraform plan failed"
        exit 1
    }

    # Apply the changes
    terraform apply tfplan

    if ($LASTEXITCODE -ne 0) {
        Write-Error "Terraform apply failed"
        exit 1
    }

    Write-Success "Infrastructure deployed successfully"
}

function Test-Integration {
    Write-Info "Running integration tests..."

    # Get the ALB DNS name
    $albDns = terraform output -raw alb_dns_name

    if ($LASTEXITCODE -ne 0 -or [string]::IsNullOrEmpty($albDns)) {
        Write-Error "Could not get ALB DNS name from Terraform output"
        return $false
    }

    Write-Info "Testing against ALB: $albDns"

    # Wait for services to be ready
    Write-Info "Waiting for services to be healthy..."
    Start-Sleep -Seconds 120

    # Test health endpoint
    try {
        $response = Invoke-WebRequest -Uri "http://$albDns/health" -TimeoutSec 30
        if ($response.StatusCode -eq 200) {
            Write-Success "Health check passed"
        } else {
            Write-Error "Health check failed with status: $($response.StatusCode)"
            return $false
        }
    } catch {
        Write-Error "Health check failed: $($_.Exception.Message)"
        return $false
    }

    # Test API health endpoint
    try {
        $response = Invoke-WebRequest -Uri "http://$albDns/api/health" -TimeoutSec 30
        if ($response.StatusCode -eq 200) {
            Write-Success "API health check passed"
        } else {
            Write-Error "API health check failed with status: $($response.StatusCode)"
            return $false
        }
    } catch {
        Write-Error "API health check failed: $($_.Exception.Message)"
        return $false
    }

    Write-Success "All integration tests passed!"
    return $true
}

function Invoke-Cleanup {
    Write-Info "Cleaning up temporary files..."

    # Remove Terraform plan
    if (Test-Path tfplan) {
        Remove-Item tfplan
    }

    # Remove local Docker images (optional)
    # docker rmi btng:$DOCKER_IMAGE_TAG 2>$null

    Write-Success "Cleanup complete"
}

function Invoke-Main {
    Write-Info "🚀 BTNG Sovereign Gold Standard - AWS Deployment"
    Write-Host "=================================================="

    Test-Dependencies
    Build-DockerImage
    Authenticate-Ecr
    New-EcrRepository
    Push-DockerImage
    Setup-Terraform
    Deploy-Infrastructure

    if (Test-Integration) {
        Invoke-Cleanup

        # Get outputs
        $albDns = terraform output -raw alb_dns_name
        $ecrUri = terraform output -raw ecr_repository_url

        Write-Host ""
        Write-Success "🎉 Deployment Complete!"
        Write-Host ""
        Write-Host "🌐 Application URL: http://$albDns"
        Write-Host "🔗 Health Check: http://$albDns/health"
        Write-Host "🔗 API Health Check: http://$albDns/api/health"
        Write-Host "🐳 ECR Repository: $ecrUri"
        Write-Host ""
        Write-Host "📋 Next Steps:"
        Write-Host "  1. Update your DNS to point to the ALB: $albDns"
        Write-Host "  2. Configure HTTPS with AWS Certificate Manager"
        Write-Host "  3. Set up monitoring and alerting"
        Write-Host "  4. Configure auto-scaling based on your needs"
    } else {
        Write-Error "Integration tests failed. Please check the deployment."
        exit 1
    }
}

# Handle command line arguments
switch ($Action) {
    "build" {
        Test-Dependencies
        Build-DockerImage
    }
    "push" {
        Test-Dependencies
        Authenticate-Ecr
        New-EcrRepository
        Push-DockerImage
    }
    "infra" {
        Setup-Terraform
        Deploy-Infrastructure
    }
    "test" {
        Test-Integration
    }
    "cleanup" {
        Invoke-Cleanup
    }
    default {
        Invoke-Main
    }
}