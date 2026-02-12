package goauth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
}

type Client struct {
	config     Config
	httpClient *http.Client
	projectID  uuid.UUID

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
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, fail.New(SDKRequestMarshalingErrorID).Trace(err.Error())
		}
	}

	url := fmt.Sprintf("%s%s", c.config.BaseURL, path)
	req, err := http.NewRequestWithContext(ctx, method, url, buf)
	if err != nil {
		return nil, fail.New(SDKUnknownErrorID).WithArgs(err.Error())
	}

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

	if resp.StatusCode >= 400 {
		return c.handleErrorResponse(resp)
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
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

func (c *Client) handleErrorResponse(resp *http.Response) error {
	var errResp apiErrorResponse
	body, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(body, &errResp)

	return fail.From(&httpStatusError{
		status: resp.StatusCode,
		apiID:  errResp.ErrorID,
		msg:    errResp.Message,
		traces: errResp.Trace,
	})
}
