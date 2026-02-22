# BTNG Sovereign Gold Standard - AWS Production Deployment

This guide provides complete instructions for deploying the BTNG sovereign gold standard platform to AWS using Docker containers and ECS Fargate.

## 🏗️ Architecture Overview

The deployment consists of:
- **Next.js API** (Port 3003): REST API with JWT authentication and gold price endpoints
- **Genesis Platform** (Port 8080): Kotlin JVM application for trading infrastructure
- **Nginx Reverse Proxy** (Port 80): Routes requests to appropriate services
- **Multi-stage Docker Image**: Unified container with all services
- **AWS ECS Fargate**: Serverless container orchestration
- **Application Load Balancer**: Load balancing and SSL termination
- **Auto-scaling**: CPU/memory-based scaling policies

### Service Routing

Nginx handles internal routing:
- `/api/*` → Next.js API (port 3003)
- `/genesis/*` → Genesis Platform (port 8080)
- `/health` → Nginx health check
- `/*` → Next.js API (default)

### Container Configuration

The Docker container includes:
- **Nginx** (Port 80): Reverse proxy and load balancer
- **Next.js API** (Internal port 3003): Gold price and authentication endpoints
- **Genesis Platform** (Internal port 8080): Trading infrastructure
- **Health Checks**: Automated monitoring of all services

## 📋 Prerequisites

### Required Tools
- Docker Desktop
- AWS CLI v2
- Terraform v1.2.0+
- PowerShell (Windows) or Bash (Linux/Mac)

### AWS Requirements
- AWS Account with appropriate permissions
- IAM user/role with these permissions:
  - `EC2FullAccess`
  - `ECSFULLAccess`
  - `ECRFullAccess`
  - `ELBFullAccess`
  - `IAMFullAccess`
  - `CloudWatchFullAccess`
  - `AutoScalingFullAccess`

### Environment Setup
1. Configure AWS CLI:
   ```bash
   aws configure
   ```

2. Set your desired AWS region (default: us-east-1):
   ```bash
   aws configure set region us-east-1
   ```

## 🚀 Quick Deployment

### Option 1: Full Automated Deployment (Recommended)

Run the complete deployment with a single command:

```powershell
# Windows PowerShell
.\deploy-aws.ps1

# Or specify a specific action
.\deploy-aws.ps1 -Action build    # Build Docker image only
.\deploy-aws.ps1 -Action push     # Push to ECR only
.\deploy-aws.ps1 -Action infra    # Deploy infrastructure only
.\deploy-aws.ps1 -Action test     # Run tests only
```

### Option 2: Manual Step-by-Step Deployment

#### 1. Configure Environment Variables

Edit `terraform.tfvars` with your specific values:

```hcl
# terraform.tfvars
aws_region = "us-east-1"
environment = "production"

# Application settings
app_name = "btng"
app_version = "1.0.0"

# Database settings (if using RDS)
db_username = "btng_admin"
db_password = "your_secure_password_here"

# JWT and API keys
jwt_secret = "your_jwt_secret_here"
api_key = "your_api_key_here"

# Gold price API settings
gold_api_key = "your_gold_api_key_here"
gold_api_url = "https://api.goldprice.com/v1"

# Email settings (for notifications)
smtp_host = "smtp.gmail.com"
smtp_port = "587"
smtp_username = "your_email@gmail.com"
smtp_password = "your_app_password"

# Monitoring and alerting
alert_email = "admin@yourcompany.com"
```

#### 2. Build and Push Docker Image

```powershell
# Build the multi-stage Docker image
docker build -t btng:latest .

# Authenticate with ECR
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin YOUR_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com

# Create ECR repository (if it doesn't exist)
aws ecr create-repository --repository-name btng --region us-east-1

# Tag and push the image
docker tag btng:latest YOUR_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/btng:latest
docker push YOUR_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/btng:latest
```

#### 3. Deploy Infrastructure

```powershell
# Initialize Terraform
terraform init

# Validate configuration
terraform validate

# Plan the deployment
terraform plan -var-file="terraform.tfvars" -out=tfplan

# Apply the changes
terraform apply tfplan
```

## 🔍 Verification and Testing

### Health Checks

After deployment, verify the services are running:

```bash
# Get the ALB DNS name
terraform output alb_dns_name

# Test health endpoints
curl http://YOUR_ALB_DNS/api/health
curl http://YOUR_ALB_DNS/api/btng/gold/price/status
```

### Integration Tests

The deployment script includes automated integration tests. You can also run them manually:

```powershell
.\deploy-aws.ps1 -Action test
```

## 📊 Monitoring and Observability

### CloudWatch Metrics

The deployment includes:
- ECS service metrics (CPU, memory utilization)
- ALB request metrics (latency, error rates)
- Application logs (API requests, errors)

### Accessing Logs

```bash
# View ECS service logs
aws logs tail /ecs/btng --follow --region us-east-1

# View ALB access logs
aws s3 ls s3://btng-alb-logs-YYYY-MM-DD/
```

### Setting up Alerts

Configure CloudWatch alarms for:
- High CPU/memory usage
- Service health check failures
- ALB 5xx errors

## 🔒 Security Configuration

### SSL/TLS Setup

1. Request an ACM certificate:
   ```bash
   aws acm request-certificate --domain-name yourdomain.com --validation-method DNS
   ```

2. Update the ALB listener in Terraform to use HTTPS

3. Configure DNS to point to the ALB

### Security Groups

The deployment creates security groups that:
- Allow inbound traffic on ports 80/443 from anywhere
- Allow outbound traffic to required services
- Restrict database access to ECS tasks only

### Secrets Management

Sensitive data is stored in:
- AWS Systems Manager Parameter Store (for non-sensitive config)
- AWS Secrets Manager (for sensitive data like database passwords)

## 🔄 Updates and Rollbacks

### Updating the Application

1. Build and push a new Docker image with a new tag
2. Update the ECS service to use the new image tag
3. Monitor the deployment with CloudWatch

### Rolling Back

```bash
# Update to previous image tag
aws ecs update-service --cluster btng-cluster --service btng-service --task-definition btng:previous-tag
```

## 🏗️ Infrastructure Components

### VPC and Networking
- VPC with public and private subnets across 2 AZs
- Internet Gateway for public access
- NAT Gateway for private subnet outbound traffic
- Route tables and security groups

### Load Balancing
- Application Load Balancer with path-based routing
- Health checks for both services
- SSL termination support

### Container Orchestration
- ECS Cluster with Fargate launch type
- Task definition with both services
- Auto-scaling policies based on CPU/memory

### Storage and Database
- ECR for container images
- (Optional) RDS for persistent data
- CloudWatch for logs and metrics

## 🐛 Troubleshooting

### Common Issues

1. **Docker Build Fails**
   - Ensure all required files are present
   - Check for syntax errors in Docker commands

2. **ECR Push Fails**
   - Verify AWS credentials and permissions
   - Check ECR repository exists

3. **Terraform Apply Fails**
   - Review error messages for specific resource issues
   - Check AWS limits and quotas

4. **Service Health Check Fails**
   - Check ECS task logs for application errors
   - Verify environment variables are set correctly

### Debugging Commands

```bash
# Check ECS service status
aws ecs describe-services --cluster btng-cluster --services btng-service

# View task logs
aws ecs describe-tasks --cluster btng-cluster --tasks TASK_ID
aws logs get-log-events --log-group-name /ecs/btng --log-stream-name ecs/btng/TASK_ID

# Check ALB target health
aws elbv2 describe-target-health --target-group-arn TARGET_GROUP_ARN
```

## 📈 Scaling and Performance

### Auto-scaling Configuration

The deployment includes auto-scaling based on:
- CPU utilization > 70%
- Memory utilization > 80%
- Minimum 2 tasks, maximum 10 tasks

### Performance Optimization

- Use appropriate instance sizes for your workload
- Configure connection pooling for databases
- Implement caching for frequently accessed data
- Monitor and optimize API response times

## 🧹 Cleanup

To remove all resources:

```powershell
# Destroy infrastructure
terraform destroy -var-file="terraform.tfvars"

# Remove ECR images
aws ecr batch-delete-image --repository-name btng --image-ids imageTag=latest

# Delete ECR repository
aws ecr delete-repository --repository-name btng --force
```

## 📞 Support

For issues or questions:
1. Check the troubleshooting section above
2. Review AWS documentation for specific services
3. Check application logs in CloudWatch
4. Contact the development team

---

**Deployment Complete!** 🎉

Your BTNG sovereign gold standard platform is now running on AWS with full production infrastructure, monitoring, and auto-scaling capabilities.