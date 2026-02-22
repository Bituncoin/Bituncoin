# BTNG Sovereign Gold Standard - Terraform Variables
# Update these values with your specific configuration

# AWS Configuration
aws_region = "us-east-1"
environment = "production"

# Application Settings
app_name = "btng"
app_version = "1.0.0"

# Database Settings (if using RDS - optional)
db_username = "btng_admin"
db_password = "your_secure_db_password_here"

# Authentication & Security
jwt_secret = "your-jwt-secret-here-at-least-256-bits"
api_secret = "your-api-secret-here"

# Gold Price API Configuration
gold_api_key = "your_gold_api_key_here"
gold_api_url = "https://api.goldprice.com/v1"

# Blockchain Configuration
eth_rpc_url = "https://eth-sepolia.g.alchemy.com/v2/YOUR_ALCHEMY_KEY"

# Email Configuration (for notifications)
smtp_host = "smtp.gmail.com"
smtp_port = "587"
smtp_username = "your_email@gmail.com"
smtp_password = "your_app_password_here"

# Monitoring & Alerting
alert_email = "admin@yourcompany.com"

# Application Ports
api_port = 3003
genesis_port = 8080

# ECS Configuration
desired_count = 2
min_capacity = 2
max_capacity = 10

# CPU and Memory for Fargate tasks
cpu = "1024"  # 1 vCPU
memory = "2048"  # 2 GB RAM