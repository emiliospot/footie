# Footie Infrastructure - Terraform

This directory contains Terraform configurations for deploying the Footie application to AWS.

## Architecture

- **VPC**: Multi-AZ VPC with public and private subnets
- **ECS**: Fargate containers for backend and frontend
- **RDS**: PostgreSQL database
- **ElastiCache**: Redis for caching
- **S3**: Static asset storage
- **CloudFront**: CDN for global content delivery
- **ALB**: Application Load Balancer
- **ECR**: Container registries

## Prerequisites

1. AWS CLI configured with appropriate credentials
2. Terraform >= 1.5.0
3. S3 bucket for Terraform state (create manually first)
4. DynamoDB table for state locking (create manually first)

## Setup

### 1. Create S3 Bucket for Terraform State

```bash
aws s3api create-bucket \
  --bucket footie-terraform-state \
  --region eu-west-1 \
  --create-bucket-configuration LocationConstraint=eu-west-1

aws s3api put-bucket-versioning \
  --bucket footie-terraform-state \
  --versioning-configuration Status=Enabled
```

### 2. Create DynamoDB Table for State Locking

```bash
aws dynamodb create-table \
  --table-name footie-terraform-locks \
  --attribute-definitions AttributeName=LockID,AttributeType=S \
  --key-schema AttributeName=LockID,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  --region eu-west-1
```

### 3. Create ECR Repositories

```bash
aws ecr create-repository --repository-name footie-backend --region eu-west-1
aws ecr create-repository --repository-name footie-frontend --region eu-west-1
```

### 4. Create terraform.tfvars

```hcl
aws_region         = "eu-west-1"
environment        = "production"
db_username        = "footie_admin"
db_password        = "CHANGE_ME_SECURE_PASSWORD"
jwt_secret         = "CHANGE_ME_SECURE_JWT_SECRET"
s3_bucket_name     = "footie-assets-production"
backend_image      = "YOUR_AWS_ACCOUNT.dkr.ecr.eu-west-1.amazonaws.com/footie-backend:latest"
frontend_image     = "YOUR_AWS_ACCOUNT.dkr.ecr.eu-west-1.amazonaws.com/footie-frontend:latest"
```

## Deployment

### Initialize Terraform

```bash
terraform init
```

### Plan Infrastructure

```bash
terraform plan
```

### Apply Infrastructure

```bash
terraform apply
```

### Destroy Infrastructure

```bash
terraform destroy
```

## Modules

- `vpc/`: VPC, subnets, NAT gateways
- `rds/`: PostgreSQL database
- `redis/`: ElastiCache Redis cluster
- `ecs/`: ECS cluster, services, and tasks
- `s3/`: S3 buckets for static assets
- `cloudfront/`: CloudFront distribution

## Outputs

After successful deployment, Terraform will output:

- VPC ID
- Database endpoint
- Redis endpoint
- Load Balancer DNS
- CloudFront domain
- S3 bucket name

## Security Notes

1. Never commit `terraform.tfvars` with sensitive data
2. Use AWS Secrets Manager for production secrets
3. Enable MFA for AWS accounts
4. Rotate credentials regularly
5. Review security groups and network ACLs

## Cost Optimization

- Use t4g instances for cost savings
- Enable auto-scaling for ECS services
- Use S3 lifecycle policies
- Consider Reserved Instances for predictable workloads
