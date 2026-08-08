package payssage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Config struct {
	BaseURL string // Payssage API base URL
	APIKey  string // IdentityX access API key resolving to the platform/service actor
	AppURL  string // Payssage app base URL (used to build provider OAuth callback URLs)
}

type Client struct {
	baseURL    string
	apiKey     string
	appURL     string
	httpClient *http.Client
}

func New(cfg Config) *Client {
	return &Client{
		baseURL:    strings.TrimRight(cfg.BaseURL, "/"),
		apiKey:     cfg.APIKey,
		appURL:     strings.TrimRight(cfg.AppURL, "/"),
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *Client) do(ctx context.Context, method, path string, body, out any) error {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("payssage: marshal request: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return fmt.Errorf("payssage: build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("payssage: http: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode >= 400 {
		var apiErr APIError
		_ = json.NewDecoder(resp.Body).Decode(&apiErr)
		apiErr.StatusCode = resp.StatusCode
		return &apiErr
	}

	if out != nil {
		raw := struct {
			Data json.RawMessage `json:"data"`
		}{}
		err = json.NewDecoder(resp.Body).Decode(&raw)
		if err != nil {
			return fmt.Errorf("payssage: decode envelope: %w", err)
		}
		err = json.Unmarshal(raw.Data, out)
		if err != nil {
			return fmt.Errorf("payssage: decode data: %w", err)
		}
	}

	return nil
}

type APIError struct {
	StatusCode int    `json:"-"`
	Module     string `json:"module"`
	Message    string `json:"message"`
	Code       int    `json:"code"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("payssage: api error %d: %s", e.StatusCode, e.Message)
}

func IsNotFound(err error) bool {
	var apiErr *APIError
	if err == nil {
		return false
	}
	if e, ok := errors.AsType[*APIError](err); ok {
		return e.StatusCode == http.StatusNotFound || errors.Is(e, apiErr)
	}
	return false
}

func IsUnauthorized(err error) bool {
	if e, ok := errors.AsType[*APIError](err); ok {
		return e.StatusCode == http.StatusUnauthorized
	}
	return false
}

func IsConflict(err error) bool {
	if e, ok := errors.AsType[*APIError](err); ok {
		return e.StatusCode == http.StatusConflict
	}
	return false
}
