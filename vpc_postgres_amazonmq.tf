terraform {
  required_version = ">= 1.5.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region

  # For production, prefer standard AWS auth (env vars / shared config / IAM role).
  access_key = var.use_localstack ? "test" : var.aws_access_key_id
  secret_key = var.use_localstack ? "test" : var.aws_secret_access_key
  token      = var.use_localstack ? null : var.aws_session_token

  skip_credentials_validation = var.use_localstack
  skip_metadata_api_check     = var.use_localstack
  skip_requesting_account_id  = var.use_localstack
  s3_use_path_style           = var.use_localstack

  endpoints {
    ec2   = var.use_localstack ? var.localstack_endpoint : null
    rds   = var.use_localstack ? var.localstack_endpoint : null
    mq    = var.use_localstack ? var.localstack_endpoint : null
    iam   = var.use_localstack ? var.localstack_endpoint : null
    sts   = var.use_localstack ? var.localstack_endpoint : null
    logs  = var.use_localstack ? var.localstack_endpoint : null
    route53 = var.use_localstack ? var.localstack_endpoint : null
  }
}

variable "project_name" {
  type    = string
  default = "core-infra"
}

variable "environment" {
  type    = string
  default = "prod"
}

variable "aws_region" {
  type    = string
  default = "us-east-1"
}

variable "use_localstack" {
  type    = bool
  default = false
}

variable "localstack_endpoint" {
  type    = string
  default = "http://localhost:4566"
}

variable "aws_access_key_id" {
  type    = string
  default = null
}

variable "aws_secret_access_key" {
  type      = string
  default   = null
  sensitive = true
}

variable "aws_session_token" {
  type      = string
  default   = null
  sensitive = true
}

variable "vpc_cidr" {
  type    = string
  default = "10.42.0.0/16"
}

variable "public_subnet_cidr" {
  type    = string
  default = "10.42.0.0/24"
}

variable "private_subnet_cidrs" {
  type    = list(string)
  default = ["10.42.1.0/24", "10.42.2.0/24"]

  validation {
    condition     = length(var.private_subnet_cidrs) >= 2
    error_message = "Provide at least 2 private subnet CIDRs (RDS subnet group requirement)."
  }
}

variable "availability_zones" {
  type    = list(string)
  default = []
}

variable "allowed_public_ingress_cidrs" {
  type    = list(string)
  default = ["0.0.0.0/0"]
}

variable "enable_postgres" {
  type    = bool
  default = true
}

variable "postgres_db_name" {
  type    = string
  default = "appdb"
}

variable "postgres_username" {
  type    = string
  default = "appuser"
}

variable "postgres_password" {
  type      = string
  sensitive = true
  default   = "change-me-now"
}

variable "postgres_engine_version" {
  type    = string
  default = "16.3"
}

variable "postgres_instance_class" {
  type    = string
  default = "db.t4g.micro"
}

variable "postgres_allocated_storage" {
  type    = number
  default = 20
}

variable "postgres_multi_az" {
  type    = bool
  default = false
}

variable "postgres_storage_encrypted" {
  type    = bool
  default = true
}

variable "postgres_port" {
  type    = number
  default = 5432
}

variable "postgres_skip_final_snapshot" {
  type    = bool
  default = true
}

variable "enable_amazonmq" {
  type    = bool
  default = true
}

variable "mq_engine_type" {
  type    = string
  default = "RabbitMQ"

  validation {
    condition     = contains(["ActiveMQ", "RabbitMQ"], var.mq_engine_type)
    error_message = "mq_engine_type must be ActiveMQ or RabbitMQ."
  }
}

variable "mq_engine_version" {
  type    = string
  default = "3.13"
}

variable "mq_host_instance_type" {
  type    = string
  default = "mq.t3.micro"
}

variable "mq_deployment_mode" {
  type    = string
  default = "SINGLE_INSTANCE"

  validation {
    condition     = contains(["SINGLE_INSTANCE", "ACTIVE_STANDBY_MULTI_AZ"], var.mq_deployment_mode)
    error_message = "mq_deployment_mode must be SINGLE_INSTANCE or ACTIVE_STANDBY_MULTI_AZ."
  }
}

variable "mq_username" {
  type    = string
  default = "mqadmin"
}

variable "mq_password" {
  type      = string
  sensitive = true
  default   = "ChangeMeNow123!"
}

variable "mq_port" {
  type    = number
  default = 5671
}

variable "tags" {
  type    = map(string)
  default = {}
}

locals {
  name_prefix = "${var.project_name}-${var.environment}"
  azs         = length(var.availability_zones) > 0 ? var.availability_zones : ["${var.aws_region}a", "${var.aws_region}b"]

  private_subnets = {
    for idx, cidr in var.private_subnet_cidrs :
    tostring(idx) => {
      cidr = cidr
      az   = local.azs[idx % length(local.azs)]
    }
  }

  merged_tags = merge(
    {
      Project     = var.project_name
      Environment = var.environment
      ManagedBy   = "Terraform"
    },
    var.tags
  )
}

resource "aws_vpc" "main" {
  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = merge(local.merged_tags, { Name = "${local.name_prefix}-vpc" })
}

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id
  tags   = merge(local.merged_tags, { Name = "${local.name_prefix}-igw" })
}

resource "aws_subnet" "public" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = var.public_subnet_cidr
  availability_zone       = local.azs[0]
  map_public_ip_on_launch = true

  tags = merge(local.merged_tags, { Name = "${local.name_prefix}-public-1" })
}

resource "aws_subnet" "private" {
  for_each = local.private_subnets

  vpc_id            = aws_vpc.main.id
  cidr_block        = each.value.cidr
  availability_zone = each.value.az

  tags = merge(local.merged_tags, { Name = "${local.name_prefix}-private-${each.key}" })
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }

  tags = merge(local.merged_tags, { Name = "${local.name_prefix}-public-rt" })
}

resource "aws_route_table_association" "public" {
  subnet_id      = aws_subnet.public.id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table" "private" {
  vpc_id = aws_vpc.main.id
  tags   = merge(local.merged_tags, { Name = "${local.name_prefix}-private-rt" })
}

resource "aws_route_table_association" "private" {
  for_each = aws_subnet.private

  subnet_id      = each.value.id
  route_table_id = aws_route_table.private.id
}

# Public endpoint SG (attach to your public app/ingress instance in this subnet).
resource "aws_security_group" "public_endpoint" {
  name        = "${local.name_prefix}-public-endpoint-sg"
  description = "Public endpoint access"
  vpc_id      = aws_vpc.main.id

  dynamic "ingress" {
    for_each = toset(var.allowed_public_ingress_cidrs)
    content {
      from_port   = 443
      to_port     = 443
      protocol    = "tcp"
      cidr_blocks = [ingress.value]
    }
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(local.merged_tags, { Name = "${local.name_prefix}-public-endpoint-sg" })
}

resource "aws_security_group" "postgres" {
  name        = "${local.name_prefix}-postgres-sg"
  description = "Private PostgreSQL access from public endpoint"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port       = var.postgres_port
    to_port         = var.postgres_port
    protocol        = "tcp"
    security_groups = [aws_security_group.public_endpoint.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(local.merged_tags, { Name = "${local.name_prefix}-postgres-sg" })
}

resource "aws_security_group" "mq" {
  name        = "${local.name_prefix}-mq-sg"
  description = "Private AmazonMQ access from public endpoint"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port       = var.mq_port
    to_port         = var.mq_port
    protocol        = "tcp"
    security_groups = [aws_security_group.public_endpoint.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(local.merged_tags, { Name = "${local.name_prefix}-mq-sg" })
}

resource "aws_db_subnet_group" "postgres" {
  count = var.enable_postgres ? 1 : 0

  name       = "${local.name_prefix}-postgres-subnets"
  subnet_ids = [for s in aws_subnet.private : s.id]

  tags = merge(local.merged_tags, { Name = "${local.name_prefix}-postgres-subnet-group" })
}

resource "aws_db_instance" "postgres" {
  count = var.enable_postgres ? 1 : 0

  identifier             = "${local.name_prefix}-postgres"
  engine                 = "postgres"
  engine_version         = var.postgres_engine_version
  instance_class         = var.postgres_instance_class
  allocated_storage      = var.postgres_allocated_storage
  db_name                = var.postgres_db_name
  username               = var.postgres_username
  password               = var.postgres_password
  port                   = var.postgres_port
  db_subnet_group_name   = aws_db_subnet_group.postgres[0].name
  vpc_security_group_ids = [aws_security_group.postgres.id]
  publicly_accessible    = false
  multi_az               = var.postgres_multi_az
  storage_encrypted      = var.use_localstack ? false : var.postgres_storage_encrypted
  skip_final_snapshot    = var.postgres_skip_final_snapshot
  apply_immediately      = true

  tags = merge(local.merged_tags, { Name = "${local.name_prefix}-postgres" })
}

resource "aws_mq_broker" "main" {
  count = var.enable_amazonmq ? 1 : 0

  broker_name         = "${local.name_prefix}-mq"
  engine_type         = var.mq_engine_type
  engine_version      = var.mq_engine_version
  host_instance_type  = var.mq_host_instance_type
  deployment_mode     = var.mq_deployment_mode
  security_groups     = [aws_security_group.mq.id]
  publicly_accessible = false
  subnet_ids = var.mq_deployment_mode == "ACTIVE_STANDBY_MULTI_AZ" ? [for s in aws_subnet.private : s.id] : [sort([for s in aws_subnet.private : s.id])[0]]
  apply_immediately = true

  user {
    username = var.mq_username
    password = var.mq_password
  }

  tags = local.merged_tags
}

output "vpc_id" {
  value = aws_vpc.main.id
}

output "public_subnet_id" {
  value = aws_subnet.public.id
}

output "private_subnet_ids" {
  value = [for s in aws_subnet.private : s.id]
}

output "public_endpoint_security_group_id" {
  value = aws_security_group.public_endpoint.id
}

output "postgres_endpoint" {
  value = try(aws_db_instance.postgres[0].address, null)
}

output "postgres_port" {
  value = var.postgres_port
}

output "amazonmq_primary_endpoint" {
  value = try(aws_mq_broker.main[0].instances[0].endpoints[0], null)
}
