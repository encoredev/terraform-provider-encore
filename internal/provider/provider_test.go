// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hasura/go-graphql-client"
)

func newTestPlatformClient(string) PlatformClient {
	tp := &TestPlatformClient{}
	tp.gql = graphql.NewClient("http://localhost:8080/graphql", tp)
	return tp
}

type TestPlatformClient struct {
	gql *graphql.Client
}

func (t TestPlatformClient) Do(req *http.Request) (*http.Response, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	reqBody := struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables"`
	}{}
	if err := json.Unmarshal(body, &reqBody); err != nil {
		return nil, err
	}
	fileStream, err := os.Open(fmt.Sprintf("testdata/%s.json", reqBody.Variables["envName"]))
	if err != nil {
		return nil, err
	}
	return &http.Response{
		StatusCode: 200,
		Body:       fileStream,
	}, nil

}

func (t TestPlatformClient) Auth(ctx context.Context, authKey string) error {
	return nil
}

func (t TestPlatformClient) Call(ctx context.Context, method, path string, reqParams, respParams interface{}) error {
	return nil
}

func (t TestPlatformClient) GQL() *graphql.Client {
	return t.gql
}

func (t TestPlatformClient) AppSlug() string {
	return "test"
}

var _ PlatformClient = &TestPlatformClient{}

// testV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"encore": providerserver.NewProtocol6WithError(newForTest("test", newTestPlatformClient)()),
}

func newForTest(version string, clientFactory func(string) PlatformClient) func() provider.Provider {
	return func() provider.Provider {
		return &EncoreProvider{
			version:       version,
			clientFactory: clientFactory,
		}
	}
}
