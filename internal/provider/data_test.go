package provider

import (
	"context"
	"os"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestData(t *testing.T) {
	c := qt.New(t)
	c.Skip("skipping test in CI")
	ctx := context.Background()
	client := NewPlatformClient("test")
	err := client.Auth(ctx, os.Getenv("ENCORE_AUTH_KEY"))
	c.Assert(err, qt.IsNil)
	nd := NewNeedsData(client, "staging", []func() datasource.DataSource{
		NewPubSubTopic,
		NewPubSubSubscription,
		NewDatabase,
		NewCache,
		NewService,
		NewGateway,
	})
	_, diags := nd.Get(ctx, "need.Topic", "", "test")
	c.Assert(diags, qt.HasLen, 0)
}
