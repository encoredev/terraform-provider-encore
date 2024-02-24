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

type GCPSQLServer struct {
	SelfLink string `tf:"id"`
	Network  GCPNetwork
	SslCert  GCPSSLCert
}

type GCPSSLCert struct {
	Fingerprint string
}

type GCPNetwork struct {
	SelfLink string `tf:"id"`
}

type AWSSQLServer struct {
	Arn            string
	VPC            AWSVPC
	SubnetGroup    AWSSubnetGroup
	SecurityGroup  AWSSecurityGroup
	ParameterGroup AWSParameterGroup
}

type AWSParameterGroup struct {
	Arn string
}

type AWSSecurityGroup struct {
	ID string
}

type AWSSubnet struct {
	Arn string
	Az  string
	Vpc AWSVPC
}

type AWSVPC struct {
	ID string
}

func (a *SQLDatabase) GetDocs() map[string]string {
	return map[string]string{
		"arn": "The [ARN](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference-arns.html) for this  sns topic",
	}
}
