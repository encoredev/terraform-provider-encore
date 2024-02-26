package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func NewDatabase() datasource.DataSource {
	return NewEncoreDataSource(
		"need.Database",
		"sql_database",
		"Encore provisioned database information",
		"SQLDatabase")
}

type SQLDatabase struct {
	SQLServer `graphql:"server"`
}

type SQLServer struct {
	AwsRds      AWSSQLServer `graphql:"... on AWSSQLServer"`
	GcpCloudSQL GCPSQLServer `graphql:"... on GCPSQLServer"`
}

func (a *SQLServer) GetDocs() map[string]string {
	return map[string]string{
		"aws_rds":       "Set if the database server instance is an AWS RDS instance",
		"gcp_cloud_sql": "Set if the database server instance is a GCP Cloud SQL instance",
	}
}

type GCPSQLServer struct {
	SelfLink string `tf:"id"`
	Network  GCPNetwork
	SslCert  GCPSSLCert
}

func (a *GCPSQLServer) GetDocs() map[string]string {
	return map[string]string{
		"id":       "The [id](https://cloud.google.com/apis/design/resource_names#relative_resource_name) in the form of `projects/{project}/instances/{instance}`",
		"network":  "The [network](https://cloud.google.com/vpc/docs/vpc) that the database instance is connected to",
		"ssl_cert": "The [SSL certificate](https://cloud.google.com/sql/docs/mysql/configure-ssl-instance) for the database instance",
	}
}

type GCPSSLCert struct {
	Fingerprint string
}

func (a *GCPSSLCert) GetDocs() map[string]string {
	return map[string]string{
		"fingerprint": "The [fingerprint](https://cloud.google.com/sql/docs/mysql/configure-ssl-instance) of the SSL certificate",
	}
}

type GCPNetwork struct {
	SelfLink string `tf:"id"`
}

func (a *GCPNetwork) GetDocs() map[string]string {
	return map[string]string{
		"id": "The [id](https://cloud.google.com/apis/design/resource_names#relative_resource_name) in the form of `projects/{project}/global/networks/{network}`",
	}
}

type AWSSQLServer struct {
	Arn            string
	VPC            AWSVPC
	SubnetGroup    AWSSubnetGroup
	SecurityGroup  AWSSecurityGroup
	ParameterGroup AWSParameterGroup
}

func (a *AWSSQLServer) GetDocs() map[string]string {
	return map[string]string{
		"arn":             "The [ARN](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference-arns.html) for the database server instance",
		"vpc":             "The [VPC](https://docs.aws.amazon.com/vpc/latest/userguide/what-is-amazon-vpc.html) that the database instance is connected to",
		"subnet_group":    "The [subnet group](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_VPC.WorkingWithRDSInstanceinaVPC.html) that the database instance is connected to",
		"security_group":  "The [security group](https://docs.aws.amazon.com/vpc/latest/userguide/VPC_SecurityGroups.html) that the database instance is connected to",
		"parameter_group": "The [parameter group](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_WorkingWithParamGroups.html) that the database instance uses",
	}
}

type AWSParameterGroup struct {
	Arn string
}

func (a *AWSParameterGroup) GetDocs() map[string]string {
	return map[string]string{
		"arn": "The [ARN](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference-arns.html) for the parameter group",
	}
}

type AWSSecurityGroup struct {
	ID string
}

func (a *AWSSecurityGroup) GetDocs() map[string]string {
	return map[string]string{
		"id": "The [id](https://docs.aws.amazon.com/vpc/latest/userguide/VPC_SecurityGroups.html) for the security group",
	}
}

type AWSSubnet struct {
	Arn string
	Az  string
	Vpc AWSVPC
}

func (a AWSSubnet) GetDocs() map[string]string {
	return map[string]string{
		"arn": "The [ARN](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference-arns.html) for the subnet",
		"az":  "The [availability zone](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html) for the subnet",
		"vpc": "The [VPC](https://docs.aws.amazon.com/vpc/latest/userguide/what-is-amazon-vpc.html) that the subnet is connected to",
	}

}

type AWSVPC struct {
	ID string
}

func (a *AWSVPC) GetDocs() map[string]string {
	return map[string]string{
		"id": "The [id](https://docs.aws.amazon.com/vpc/latest/userguide/what-is-amazon-vpc.html) for the VPC",
	}

}
