package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExampleDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccExampleDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.encore_pubsub_topic.test", "env", "@primary"),
				),
			},
		},
	})
}

const testAccExampleDataSourceConfig = `
provider "encore" {}

data "encore_pubsub_topic" "test" {
    name = "ordered"
	env = "gke"
}
`
