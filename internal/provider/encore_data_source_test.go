package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAWSSubnets(res, prefix string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(res, prefix+".subnets.0.arn", "arn:aws:ec2:region:account:subnet/subnet"),
		resource.TestCheckResourceAttr(res, prefix+".subnets.0.az", "us-east-1"),
		resource.TestCheckResourceAttr(res, prefix+".subnets.1.arn", "arn:aws:ec2:region:account:subnet/subnet"),
		resource.TestCheckResourceAttr(res, prefix+".subnets.1.az", "us-east-1"),
	)
}

func testStepForEnv(env, cfg string, fns ...resource.TestCheckFunc) resource.TestStep {
	return resource.TestStep{
		Config: fmt.Sprintf(cfg, env),
		Check:  resource.ComposeAggregateTestCheckFunc(fns...),
	}
}
