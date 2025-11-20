output "vpc_id" {
  description = "VPC ID"
  value       = module.vpc.vpc_id
}

output "db_endpoint" {
  description = "RDS endpoint"
  value       = module.rds.db_endpoint
}

output "redis_endpoint" {
  description = "Redis endpoint"
  value       = module.redis.redis_endpoint
}

output "alb_dns_name" {
  description = "Application Load Balancer DNS name"
  value       = module.ecs.alb_dns_name
}

output "cloudfront_domain" {
  description = "CloudFront distribution domain name"
  value       = module.cloudfront.cloudfront_domain
}

output "s3_bucket_name" {
  description = "S3 bucket name"
  value       = module.s3.bucket_name
}

