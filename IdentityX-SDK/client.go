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

	Users  *UserService
	Tokens *TokenService
}

func NewClient(config Config) (*Client, error) {
	if config.BaseURL == "" {
		return nil, SdkError{"BaseURL is required", nil}
	}
	if config.APIKey == "" {
		return nil, SdkError{"APIKey is required", nil}
	}
	if config.ProjectID == uuid.Nil {
		return nil, SdkError{"ProjectID is required", nil}
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

	c.Users = &UserService{client: c}
	c.Tokens = &TokenService{client: c, cacheTTL: time.Hour}

	return c, nil
}

func (c *Client) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	var buf io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, SdkError{Message: "request marshaling error: ", Cause: err}
		}
		buf = bytes.NewReader(b)
	}

	u, err := url.JoinPath(c.config.BaseURL, path)
	if err != nil {
		return nil, SdkError{Message: "error building URL: ", Cause: err}
	}

	req, err := http.NewRequestWithContext(ctx, method, u, buf)
	if err != nil {
		return nil, SdkError{Message: "error building request: ", Cause: err}
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
		return SdkError{Message: "error executing request: ", Cause: err}
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			if c.debug {
				log.Printf("error closing response body: %v", err)
			}
		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)

	if c.debug {
		log.Printf("[DEBUG] %s %s -> %d body: %s", req.Method, req.URL.String(), resp.StatusCode, string(body))
	}

	if resp.StatusCode >= 400 {
		return c.handleErrorResponse(body)
	}

	if v != nil {
		err = json.Unmarshal(body, v)
		if err != nil {
			return SdkError{Message: "error unmarshaling response: ", Cause: err}
		}
	}

	return nil
}

func (c *Client) handleErrorResponse(body []byte) error {
	var errResp ApiError
	_ = json.Unmarshal(body, &errResp)

	return errResp
}
