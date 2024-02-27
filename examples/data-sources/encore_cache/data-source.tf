data "encore_cache" "cache" {
  name = "my-cache"
  env  = "my-env"
}

output "aws_redis" {
  value = {
    "arn" : data.encore_cache.cache.aws_redis.arn,
    "vpc" : data.encore_cache.cache.aws_redis.vpc.id,
    "security_group" : data.encore_cache.cache.aws_redis.security_group.id,
    "parameter_group" : data.encore_cache.cache.aws_redis.parameter_group.arn,
    "subnet_group" : data.encore_cache.cache.aws_redis.subnet_group.arn,
    "subnet" : data.encore_cache.cache.aws_redis.subnet_group.subnets.0.id,
    "subnet-az" : data.encore_cache.cache.aws_redis.subnet_group.subnets.0.az
  }
}

output "gcp_redis" {
  value = {
    "id" : data.encore_cache.cache.gcp_redis.id,
    "network" : data.encore_cache.cache.gcp_redis.network.id
  }
}