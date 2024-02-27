data "encore_sql_database" "database" {
  name = "my-database"
  env  = "my-env"
}

output "aws_rds" {
  value = {
    "db_name" : data.encore_sql_database.database.database_name,
    "rds_arn" : data.encore_sql_database.database.aws_rds.arn,
    "vpc" : data.encore_sql_database.database.aws_rds.vpc.id,
    "subnet_group" : data.encore_sql_database.database.aws_rds.subnet_group.arn,
    "security_group" : data.encore_sql_database.database.aws_rds.security_group.id,
    "parameter_group" : data.encore_sql_database.database.aws_rds.parameter_group.arn,
    "subnet_group" : data.encore_sql_database.database.aws_rds.subnet_group.arn,
    "subnet" : data.encore_sql_database.database.aws_rds.subnet_group.subnets.0.arn,
    "subnet_az" : data.encore_sql_database.database.aws_rds.subnet_group.subnets.0.az
  }
}

output "gcp_cloud_sql" {
  value = {
    "database_name" : data.encore_sql_database.database.database_name,
    "cloud_sql_id" : data.encore_sql_database.database.gcp_cloud_sql.id,
    "network" : data.encore_sql_database.database.gcp_cloud_sql.network.id,
    "ssl_cert_fingerprint" : data.encore_sql_database.database.gcp_cloud_sql.ssl_cert.fingerprint
  }
}
