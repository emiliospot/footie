# 🚀 Deployment Guide

This guide covers deploying the Footie application to AWS using Terraform and GitHub Actions CI/CD.

## 📋 Table of Contents

- [Prerequisites](#prerequisites)
- [Infrastructure Setup](#infrastructure-setup)
- [Environment Configuration](#environment-configuration)
- [Deployment Workflow](#deployment-workflow)
- [Manual Deployment](#manual-deployment)
- [Monitoring & Logging](#monitoring--logging)
- [Troubleshooting](#troubleshooting)

## 🔧 Prerequisites

### Required Tools

```bash
# Install Terraform
# macOS
brew install terraform

# Linux
wget https://releases.hashicorp.com/terraform/1.6.0/terraform_1.6.0_linux_amd64.zip
unzip terraform_1.6.0_linux_amd64.zip
sudo mv terraform /usr/local/bin/

# Verify installation
terraform --version
```

### AWS Credentials

1. **Create AWS Access Keys**:
   - Log in to AWS Console
   - Navigate to IAM → Users → Your User → Security Credentials
   - Create Access Key

2. **Configure AWS CLI**:

```bash
aws configure
# AWS Access Key ID: YOUR_ACCESS_KEY
# AWS Secret Access Key: YOUR_SECRET_KEY
# Default region: us-east-1
# Default output format: json
```

3. **Configure GitHub Secrets**:
   - Navigate to your repository → Settings → Secrets and variables → Actions
   - Add the following secrets:
     - `AWS_ACCESS_KEY_ID`
     - `AWS_SECRET_ACCESS_KEY`
     - `AWS_REGION` (e.g., us-east-1)
     - `DATABASE_PASSWORD` (production DB password)
     - `JWT_SECRET` (production JWT secret)
     - `REDIS_PASSWORD` (production Redis password)

## 🏗️ Infrastructure Setup

### 1. Initialize Terraform

```bash
cd infra/terraform

# Initialize Terraform
terraform init

# Validate configuration
terraform validate

# Plan infrastructure changes
terraform plan -out=tfplan

# Apply infrastructure
terraform apply tfplan
```

### 2. Terraform Modules

The infrastructure is organized into modules:

```
infra/terraform/
├── main.tf           # Main configuration
├── variables.tf      # Input variables
├── outputs.tf        # Output values
├── modules/
│   ├── vpc/          # VPC, subnets, NAT gateway
│   ├── ecs/          # ECS cluster, services, tasks
│   ├── rds/          # PostgreSQL RDS instance
│   ├── elasticache/  # Redis ElastiCache
│   ├── alb/          # Application Load Balancer
│   ├── s3/           # S3 bucket for frontend
│   └── cloudfront/   # CloudFront distribution
```

### 3. Infrastructure Components

| Component             | Purpose                 | Configuration              |
| --------------------- | ----------------------- | -------------------------- |
| **VPC**               | Network isolation       | CIDR: 10.0.0.0/16          |
| **ECS Fargate**       | Container orchestration | Auto-scaling 1-10 tasks    |
| **RDS PostgreSQL**    | Production database     | Multi-AZ, encrypted        |
| **ElastiCache Redis** | Caching layer           | Cluster mode disabled      |
| **ALB**               | Load balancing          | HTTPS with ACM certificate |
| **S3 + CloudFront**   | Frontend hosting        | Edge caching, HTTPS        |

## ⚙️ Environment Configuration

### Backend Environment Variables

Create `apps/api/.env.production`:

```bash
# Application
APP_ENV=production
APP_NAME=footie
APP_VERSION=1.0.0

# Backend API
API_PORT=8088
API_HOST=0.0.0.0

# Database (use RDS endpoint from Terraform output)
DATABASE_HOST=<RDS_ENDPOINT>
DATABASE_PORT=5432
DATABASE_NAME=footie
DATABASE_USER=footie_admin
DATABASE_PASSWORD=<SECURE_PASSWORD>
DATABASE_SSL_MODE=require

# Redis (use ElastiCache endpoint from Terraform output)
REDIS_HOST=<ELASTICACHE_ENDPOINT>
REDIS_PORT=6379
REDIS_PASSWORD=<SECURE_PASSWORD>
REDIS_DB=0

# JWT
JWT_SECRET=<SECURE_JWT_SECRET>
JWT_EXPIRY_HOURS=24
JWT_REFRESH_EXPIRY_HOURS=168

# CORS
CORS_ORIGINS=https://yourdomain.com
CORS_ALLOW_CREDENTIALS=true

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

### Frontend Environment Variables

Create `apps/web/src/environments/environment.prod.ts`:

```typescript
export const environment = {
  production: true,
  apiUrl: "https://api.yourdomain.com/api/v1",
  appName: "Footie Analytics",
};
```

## 🔄 Deployment Workflow

### Automated Deployment (Recommended)

The application automatically deploys on push to `main` branch:

```yaml
# .github/workflows/deploy.yml
name: Deploy to AWS

on:
  push:
    branches: [main]
  workflow_dispatch:

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - Build Docker images
      - Push to ECR
      - Update ECS service
      - Deploy frontend to S3
      - Invalidate CloudFront cache
```

**Deployment Steps**:

1. **Commit and push to main**:

```bash
git add .
git commit -m "feat: new feature"
git push origin main
```

2. **Monitor deployment**:
   - GitHub Actions → View workflow run
   - AWS ECS → Monitor task status
   - CloudWatch → Check logs

3. **Verify deployment**:

```bash
# Test backend
curl https://api.yourdomain.com/health

# Test frontend
curl https://yourdomain.com
```

### Manual Deployment

#### Backend (ECS Fargate)

```bash
# Build and push Docker image
cd apps/api
docker build -t footie-api:latest .

# Tag for ECR
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <AWS_ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com
docker tag footie-api:latest <AWS_ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com/footie-api:latest
docker push <AWS_ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com/footie-api:latest

# Update ECS service
aws ecs update-service \
  --cluster footie-cluster \
  --service footie-api-service \
  --force-new-deployment \
  --region us-east-1
```

#### Frontend (S3 + CloudFront)

```bash
# Build frontend
npm run build:web

# Deploy to S3
aws s3 sync dist/web/ s3://footie-frontend-bucket/ --delete

# Invalidate CloudFront cache
aws cloudfront create-invalidation \
  --distribution-id <CLOUDFRONT_DIST_ID> \
  --paths "/*"
```

## 📊 Monitoring & Logging

### CloudWatch Logs

```bash
# View backend logs
aws logs tail /ecs/footie-api --follow

# View specific log stream
aws logs get-log-events \
  --log-group-name /ecs/footie-api \
  --log-stream-name <STREAM_NAME>
```

### CloudWatch Metrics

Monitor key metrics in AWS Console:

- **ECS**: CPU/Memory utilization, task count
- **RDS**: Connections, read/write latency
- **ElastiCache**: Cache hit rate, evictions
- **ALB**: Request count, target response time

### CloudWatch Alarms

Set up alarms for:

- High CPU utilization (> 80%)
- High memory utilization (> 80%)
- Database connection errors
- 5xx error rate (> 1%)
- Target response time (> 1s)

### X-Ray Tracing (Optional)

Enable X-Ray for distributed tracing:

```go
// Add to main.go
import "github.com/aws/aws-xray-sdk-go/xray"

func main() {
    // Instrument HTTP handlers
    http.Handle("/", xray.Handler(xray.NewFixedSegmentNamer("footie-api"), router))
}
```

## 🔐 Security Best Practices

### 1. Secrets Management

Use AWS Secrets Manager or Parameter Store:

```bash
# Store secrets
aws secretsmanager create-secret \
  --name footie/production/db-password \
  --secret-string "YOUR_SECURE_PASSWORD"

# Reference in ECS task definition
{
  "secrets": [
    {
      "name": "DATABASE_PASSWORD",
      "valueFrom": "arn:aws:secretsmanager:region:account-id:secret:footie/production/db-password"
    }
  ]
}
```

### 2. IAM Roles

Follow least privilege principle:

- ECS Task Execution Role: ECR, CloudWatch Logs, Secrets Manager
- ECS Task Role: S3, RDS, ElastiCache (if needed)

### 3. Network Security

- Use private subnets for ECS tasks and RDS
- Restrict security group rules
- Enable VPC Flow Logs
- Use AWS WAF for ALB

### 4. SSL/TLS

- Use ACM for SSL certificates
- Enforce HTTPS on ALB
- Enable S3 bucket encryption
- Use encrypted RDS storage

## 🐛 Troubleshooting

### ECS Task Fails to Start

```bash
# Check task status
aws ecs describe-tasks \
  --cluster footie-cluster \
  --tasks <TASK_ARN>

# Common issues:
# - Image pull errors → Check ECR permissions
# - Resource limits → Increase CPU/memory
# - Environment variables → Verify secrets
```

### Database Connection Issues

```bash
# Test connectivity from ECS task
aws ecs execute-command \
  --cluster footie-cluster \
  --task <TASK_ID> \
  --container api \
  --command "/bin/sh" \
  --interactive

# Inside container:
psql -h <RDS_ENDPOINT> -U footie_admin -d footie
```

### High Response Times

1. **Check RDS performance**:
   - CloudWatch → RDS → DatabaseConnections
   - Enable Enhanced Monitoring

2. **Check ElastiCache hit rate**:
   - CloudWatch → ElastiCache → CacheHitRate
   - Scale up if needed

3. **Check ECS resources**:
   - CloudWatch → ECS → CPUUtilization
   - Scale horizontally if needed

### Frontend Not Loading

```bash
# Check S3 bucket policy
aws s3api get-bucket-policy --bucket footie-frontend-bucket

# Check CloudFront distribution
aws cloudfront get-distribution --id <DIST_ID>

# Test S3 directly
curl https://footie-frontend-bucket.s3.amazonaws.com/index.html
```

## 📝 Rollback Procedure

### Backend Rollback

```bash
# List task definitions
aws ecs list-task-definitions --family-prefix footie-api

# Update service to previous version
aws ecs update-service \
  --cluster footie-cluster \
  --service footie-api-service \
  --task-definition footie-api:PREVIOUS_REVISION
```

### Frontend Rollback

```bash
# Restore previous S3 version (if versioning enabled)
aws s3api list-object-versions --bucket footie-frontend-bucket

# Or redeploy previous build
aws s3 sync <previous-dist-folder>/ s3://footie-frontend-bucket/ --delete
```

## 🎉 Post-Deployment

1. **Smoke Tests**: Run automated tests against production
2. **Monitor**: Watch CloudWatch dashboards for anomalies
3. **User Testing**: Verify critical user flows
4. **Documentation**: Update CHANGELOG.md

## 📚 Additional Resources

- [AWS ECS Documentation](https://docs.aws.amazon.com/ecs/)
- [Terraform AWS Provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [AWS Well-Architected Framework](https://aws.amazon.com/architecture/well-architected/)

## 🆘 Support

For deployment issues, contact:

- DevOps Team: devops@yourdomain.com
- AWS Support: https://console.aws.amazon.com/support/

---

**Happy Deploying!** 🚀
