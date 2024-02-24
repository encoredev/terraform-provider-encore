package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/hasura/go-graphql-client"
	"golang.org/x/oauth2"
)

const DefaultBaseURL = "https://api.encore.dev"

type Error struct {
	HTTPStatus string `json:"-"`
	HTTPCode   int    `json:"-"`
	Code       string
	Detail     json.RawMessage
}

func (e Error) Error() string {
	if len(e.Detail) > 0 {
		return fmt.Sprintf("http %s: code=%s detail=%s", e.HTTPStatus, e.Code, e.Detail)
	}
	return fmt.Sprintf("http %s: code=%s", e.HTTPStatus, e.Code)
}

type OAuthData struct {
	Token   *oauth2.Token `json:"token"`
	Actor   string        `json:"actor,omitempty"` // The ID of the user or app that authorized the token.
	Email   string        `json:"email"`           // empty if logging in as an app
	AppSlug string        `json:"app_slug"`        // empty if logging in as a user
}

func NewPlatformClient(version string) *PlatformClient {
	baseURL := os.Getenv("ENCORE_API_URL")
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	return &PlatformClient{
		baseURL: baseURL,
		version: version,
		http:    http.DefaultClient,
		gql:     graphql.NewClient(baseURL+"/graphql", http.DefaultClient),
	}
}

type PlatformClient struct {
	baseURL string
	version string
	appSlug string
	http    *http.Client
	gql     *graphql.Client
}

func (p *PlatformClient) Auth(ctx context.Context, authKey string) error {
	var data OAuthData
	err := p.Call(ctx, "POST", "/login/auth-key", struct {
		AuthKey string `json:"auth_key"`
	}{authKey}, &data)
	if err != nil {
		return err
	}
	cfg := oauth2.Config{
		Endpoint: oauth2.Endpoint{
			TokenURL: p.baseURL + "/login/oauth:refresh-token",
		},
	}
	p.appSlug = data.AppSlug
	p.http = oauth2.NewClient(ctx, cfg.TokenSource(ctx, data.Token))
	p.gql = graphql.NewClient(p.baseURL+"/graphql", p.http).WithDebug(true)
	return nil
}

// Call makes a call to the API endpoint given by method and path.
// If reqParams and respParams are non-nil they are JSON-marshalled/unmarshalled.
func (p *PlatformClient) Call(ctx context.Context, method, path string, reqParams, respParams interface{}) (err error) {
	var body io.Reader
	if reqParams != nil {
		reqData, err := json.Marshal(reqParams)
		if err != nil {
			return fmt.Errorf("marshal request: %v", err)
		}
		body = bytes.NewReader(reqData)
	}

	req, err := http.NewRequestWithContext(ctx, method, p.baseURL+path, body)
	if err != nil {
		return err
	}
	if reqParams != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := p.Do(req)
	if err != nil {
		return err
	}
	return p.parsePlatformResp(resp, respParams)
}

func (p *PlatformClient) parsePlatformResp(resp *http.Response, respParams interface{}) error {
	defer func() {
		_ = resp.Body.Close()
	}()
	var respStruct struct {
		OK    bool
		Error Error
		Data  json.RawMessage
	}
	if err := json.NewDecoder(resp.Body).Decode(&respStruct); err != nil {
		return fmt.Errorf("decode response: %v", err)
	} else if !respStruct.OK {
		e := respStruct.Error
		e.HTTPCode = resp.StatusCode
		e.HTTPStatus = resp.Status
		return e
	}

	if respParams != nil {
		if err := json.Unmarshal([]byte(respStruct.Data), respParams); err != nil {
			return fmt.Errorf("decode response data: %v", err)
		}
	}
	return nil
}

func (p *PlatformClient) Do(req *http.Request) (*http.Response, error) {
	// Add a very limited amount of information for diagnostics
	req.Header.Set("User-Agent", "EncoreTF/"+p.version)
	req.Header.Set("X-Encore-Version", p.version)
	req.Header.Set("X-Encore-GOOS", runtime.GOOS)
	req.Header.Set("X-Encore-GOARCH", runtime.GOARCH)
	return p.http.Do(req)
}
