terraform {
  required_version = ">= 1.5.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    bucket         = "footie-terraform-state"
    key            = "footie/terraform.tfstate"
    region         = "eu-west-1"
    encrypt        = true
    dynamodb_table = "footie-terraform-locks"
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = "Footie"
      Environment = var.environment
      ManagedBy   = "Terraform"
    }
  }
}

# VPC
module "vpc" {
  source = "./modules/vpc"

  environment         = var.environment
  vpc_cidr            = var.vpc_cidr
  availability_zones  = var.availability_zones
}

# RDS PostgreSQL
module "rds" {
  source = "./modules/rds"

  environment         = var.environment
  vpc_id              = module.vpc.vpc_id
  private_subnet_ids  = module.vpc.private_subnet_ids
  db_name             = var.db_name
  db_username         = var.db_username
  db_password         = var.db_password
  db_instance_class   = var.db_instance_class
}

# ElastiCache Redis
module "redis" {
  source = "./modules/redis"

  environment         = var.environment
  vpc_id              = module.vpc.vpc_id
  private_subnet_ids  = module.vpc.private_subnet_ids
  node_type           = var.redis_node_type
}

# ECS Cluster
module "ecs" {
  source = "./modules/ecs"

  environment         = var.environment
  vpc_id              = module.vpc.vpc_id
  public_subnet_ids   = module.vpc.public_subnet_ids
  private_subnet_ids  = module.vpc.private_subnet_ids
  
  backend_image       = var.backend_image
  frontend_image      = var.frontend_image
  
  db_host             = module.rds.db_endpoint
  db_name             = var.db_name
  db_username         = var.db_username
  db_password         = var.db_password
  
  redis_host          = module.redis.redis_endpoint
  
  jwt_secret          = var.jwt_secret
}

# S3 for static assets
module "s3" {
  source = "./modules/s3"

  environment = var.environment
  bucket_name = var.s3_bucket_name
}

# CloudFront CDN
module "cloudfront" {
  source = "./modules/cloudfront"

  environment    = var.environment
  s3_bucket_id   = module.s3.bucket_id
  s3_bucket_domain_name = module.s3.bucket_domain_name
  alb_dns_name   = module.ecs.alb_dns_name
}

