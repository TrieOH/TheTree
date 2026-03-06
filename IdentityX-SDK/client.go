package goauth

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
)

type Config struct {
	BaseURL    string
	APIKey     string
	ProjectID  uuid.UUID
	HTTPClient *http.Client
	Debug      bool
}

type Client struct {
	config     Config
	httpClient *http.Client
	projectID  uuid.UUID
	debug      bool

	Roles       *RoleService
	Scopes      *ScopeService
	Permissions *PermissionService
	Authz       *AuthzService
	Users       *UserService
	Tokens      *TokenService
}

func NewClient(config Config) (*Client, error) {
	if config.BaseURL == "" {
		return nil, fail.New(SDKUnknownErrorID).WithArgs("BaseURL is required")
	}
	if config.APIKey == "" {
		return nil, fail.New(SDKUnknownErrorID).WithArgs("APIKey is required")
	}
	if config.ProjectID == uuid.Nil {
		return nil, fail.New(SDKUnknownErrorID).WithArgs("ProjectID is required")
	}
	config.BaseURL = strings.TrimSuffix(config.BaseURL, "/")

	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	c := &Client{
		config:     config,
		httpClient: config.HTTPClient,
		projectID:  config.ProjectID,
		debug:      config.Debug,
	}

	c.Roles = &RoleService{client: c}
	c.Scopes = &ScopeService{client: c}
	c.Permissions = &PermissionService{client: c}
	c.Authz = &AuthzService{client: c}
	c.Users = &UserService{client: c}
	c.Tokens = &TokenService{client: c}

	return c, nil
}

func (c *Client) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	var buf io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fail.New(SDKRequestMarshalingErrorID).WithArgs(err.Error())
		}
		buf = bytes.NewReader(b)
	}

	u, err := url.JoinPath(c.config.BaseURL, path)
	if err != nil {
		return nil, fail.New(SDKUnknownErrorID).WithArgs(err.Error())
	}

	req, err := http.NewRequestWithContext(ctx, method, u, buf)
	if err != nil {
		return nil, fail.New(SDKUnknownErrorID).WithArgs(err.Error())
	}

	req.Header.Set("Accept", "application/json")

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.config.APIKey != "" {
		req.Header.Set("X-API-Key", c.config.APIKey)
	}

	return req, nil
}

func (c *Client) do(req *http.Request, v any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fail.New(SDKNetworkErrorID).WithArgs(err.Error())
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if c.debug {
		log.Printf("[DEBUG] %s %s -> %d body: %s", req.Method, req.URL.String(), resp.StatusCode, string(body))
	}

	if resp.StatusCode >= 400 {
		return c.handleErrorResponse(resp, body)
	}

	if v != nil {
		err = json.Unmarshal(body, v)
		if err != nil {
			return fail.New(SDKResponseUnmarshalingErrorID).WithArgs(err.Error())
		}
	}

	return nil
}

type apiErrorResponse struct {
	ErrorID string   `json:"error_id"`
	Message string   `json:"message"`
	Trace   []string `json:"trace"`
	Code    int      `json:"code"`
}

func (c *Client) handleErrorResponse(resp *http.Response, body []byte) error {
	var errResp apiErrorResponse
	_ = json.Unmarshal(body, &errResp)

	return fail.From(&httpStatusError{
		status: resp.StatusCode,
		apiID:  errResp.ErrorID,
		msg:    errResp.Message,
		traces: errResp.Trace,
	})
}
