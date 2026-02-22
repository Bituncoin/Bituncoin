terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
  required_version = ">= 1.2.0"

  # Optional: Add backend configuration for state management
  # backend "s3" {
  #   bucket = "btng-terraform-state"
  #   key    = "btng-infrastructure.tfstate"
  #   region = "us-east-1"
  # }
}

provider "aws" {
  region = "us-east-1"

  # Optional: Add default tags
  default_tags {
    tags = {
      Project     = "BTNG"
      Environment = "production"
      ManagedBy   = "Terraform"
    }
  }
}

# ECR Repository for Docker images
resource "aws_ecr_repository" "btng" {
  name                 = "btng"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Name = "btng-ecr-repo"
  }
}

# VPC Configuration
resource "aws_vpc" "btng" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = "btng-vpc"
  }
}

resource "aws_subnet" "public_1" {
  vpc_id                  = aws_vpc.btng.id
  cidr_block              = "10.0.1.0/24"
  availability_zone       = "us-east-1a"
  map_public_ip_on_launch = true

  tags = {
    Name = "btng-public-1"
  }
}

resource "aws_subnet" "public_2" {
  vpc_id                  = aws_vpc.btng.id
  cidr_block              = "10.0.2.0/24"
  availability_zone       = "us-east-1b"
  map_public_ip_on_launch = true

  tags = {
    Name = "btng-public-2"
  }
}

# Internet Gateway
resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.btng.id

  tags = {
    Name = "btng-igw"
  }
}

# Route Table
resource "aws_route_table" "public_rt" {
  vpc_id = aws_vpc.btng.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
  }

  tags = {
    Name = "btng-public-rt"
  }
}

resource "aws_route_table_association" "public_assoc_1" {
  subnet_id      = aws_subnet.public_1.id
  route_table_id = aws_route_table.public_rt.id
}

resource "aws_route_table_association" "public_assoc_2" {
  subnet_id      = aws_subnet.public_2.id
  route_table_id = aws_route_table.public_rt.id
}

# Security Groups
resource "aws_security_group" "btng_alb_sg" {
  name_prefix = "btng-alb-sg"
  vpc_id      = aws_vpc.btng.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "HTTP"
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "HTTPS"
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "btng-alb-sg"
  }
}

resource "aws_security_group" "btng_ecs_sg" {
  name_prefix = "btng-ecs-sg"
  vpc_id      = aws_vpc.btng.id

  ingress {
    from_port       = 80
    to_port         = 80
    protocol        = "tcp"
    security_groups = [aws_security_group.btng_alb_sg.id]
    description     = "HTTP from ALB"
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "btng-ecs-sg"
  }
}

# Application Load Balancer
resource "aws_lb" "btng_alb" {
  name               = "btng-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.btng_alb_sg.id]
  subnets            = [aws_subnet.public_1.id, aws_subnet.public_2.id]

  enable_deletion_protection = false

  tags = {
    Name = "btng-alb"
  }
}

# Target Groups
resource "aws_lb_target_group" "btng" {
  name     = "btng-main"
  port     = 80
  protocol = "HTTP"
  vpc_id   = aws_vpc.btng.id

  health_check {
    enabled             = true
    path                = "/health"
    port                = "traffic-port"
    protocol            = "HTTP"
    interval            = 30
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 2
    matcher             = "200"
  }

  tags = {
    Name = "btng-main-tg"
  }
}

# ALB Listeners
resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.btng_alb.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.btng.arn
  }
}

# CloudWatch Log Group
resource "aws_cloudwatch_log_group" "btng" {
  name              = "/ecs/btng"
  retention_in_days = 30

  tags = {
    Name = "btng-log-group"
  }
}

# IAM Roles
resource "aws_iam_role" "ecs_exec" {
  name = "btng-ecs-exec-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = ""
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })

  tags = {
    Name = "btng-ecs-exec-role"
  }
}

resource "aws_iam_role_policy_attachment" "ecs_exec_AmazonECSTaskExecutionRolePolicy" {
  role       = aws_iam_role.ecs_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# ECS Cluster
resource "aws_ecs_cluster" "btng_cluster" {
  name = "btng-cluster"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }

  tags = {
    Name = "btng-ecs-cluster"
  }
}

# ECS Task Definition
resource "aws_ecs_task_definition" "btng_task" {
  family                   = "btng-task"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "1024"  # 1 vCPU
  memory                   = "2048"  # 2 GB
  execution_role_arn       = aws_iam_role.ecs_exec.arn
  task_role_arn            = aws_iam_role.ecs_exec.arn

  container_definitions = jsonencode([
    {
      name      = "btng"
      image     = "${aws_ecr_repository.btng.repository_url}:latest"
      essential = true

      portMappings = [
        {
          containerPort = 80
          hostPort      = 80
          protocol      = "tcp"
        }
      ]

      environment = [
        {
          name  = "NODE_ENV"
          value = "production"
        },
        {
          name  = "JWT_SECRET"
          value = var.jwt_secret
        },
        {
          name  = "API_SECRET"
          value = var.api_secret
        },
        {
          name  = "ETH_RPC_URL"
          value = var.eth_rpc_url
        }
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-region"        = "us-east-1"
          "awslogs-group"         = "/ecs/btng"
          "awslogs-stream-prefix" = "btng"
        }
      }

      healthCheck = {
        command = [
          "CMD-SHELL",
          "wget --no-verbose --tries=1 --spider http://localhost/health"
        ]
        interval    = 30
        timeout     = 5
        retries     = 3
        startPeriod = 60
      }
    }
  ])

  tags = {
    Name = "btng-task-definition"
  }
}

# ECS Service
resource "aws_ecs_service" "btng_service" {
  name            = "btng-service"
  cluster         = aws_ecs_cluster.btng_cluster.id
  task_definition = aws_ecs_task_definition.btng_task.arn
  desired_count   = 2

  network_configuration {
    subnets         = [aws_subnet.public_1.id, aws_subnet.public_2.id]
    security_groups = [aws_security_group.btng_ecs_sg.id]
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.btng.arn
    container_name   = "btng"
    container_port   = 80
  }

  depends_on = [
    aws_lb_listener.http
  ]

  lifecycle {
    ignore_changes = [desired_count]
  }

  tags = {
    Name = "btng-ecs-service"
  }
}

# Auto Scaling
resource "aws_appautoscaling_target" "api_scaling" {
  max_capacity       = 10
  min_capacity       = 2
  resource_id        = "service/${aws_ecs_cluster.btng_cluster.name}/${aws_ecs_service.btng_service.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  service_namespace  = "ecs"
}

resource "aws_appautoscaling_policy" "cpu_scale" {
  name               = "btng-cpu-scaling"
  policy_type        = "TargetTrackingScaling"
  resource_id        = aws_appautoscaling_target.api_scaling.resource_id
  scalable_dimension = aws_appautoscaling_target.api_scaling.scalable_dimension
  service_namespace  = aws_appautoscaling_target.api_scaling.service_namespace

  target_tracking_scaling_policy_configuration {
    predefined_metric_specification {
      predefined_metric_type = "ECSServiceAverageCPUUtilization"
    }
    target_value = 50.0
  }
}

resource "aws_appautoscaling_policy" "memory_scale" {
  name               = "btng-memory-scaling"
  policy_type        = "TargetTrackingScaling"
  resource_id        = aws_appautoscaling_target.api_scaling.resource_id
  scalable_dimension = aws_appautoscaling_target.api_scaling.scalable_dimension
  service_namespace  = aws_appautoscaling_target.api_scaling.service_namespace

  target_tracking_scaling_policy_configuration {
    predefined_metric_specification {
      predefined_metric_type = "ECSServiceAverageMemoryUtilization"
    }
    target_value = 70.0
  }
}

# Variables
variable "jwt_secret" {
  description = "JWT secret for API authentication"
  type        = string
  sensitive   = true
}

variable "api_secret" {
  description = "API secret for additional security"
  type        = string
  sensitive   = true
}

variable "eth_rpc_url" {
  description = "Ethereum RPC URL for blockchain integration"
  type        = string
  sensitive   = true
}

# Outputs
output "alb_dns_name" {
  description = "DNS name of the Application Load Balancer"
  value       = aws_lb.btng_alb.dns_name
}

output "ecr_repository_url" {
  description = "URL of the ECR repository"
  value       = aws_ecr_repository.btng.repository_url
}

output "ecs_cluster_name" {
  description = "Name of the ECS cluster"
  value       = aws_ecs_cluster.btng_cluster.name
}

output "ecs_service_name" {
  description = "Name of the ECS service"
  value       = aws_ecs_service.btng_service.name
}