package paymentsSDK

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

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func New(baseURL, apiKey string) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

type envelope[T any] struct {
	Data T `json:"data"`
}

func (c *Client) do(ctx context.Context, method, path string, body, out any) error {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("triepayments: marshal request: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return fmt.Errorf("triepayments: build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("triepayments: http: %w", err)
	}
	defer resp.Body.Close()

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
		if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
			return fmt.Errorf("triepayments: decode envelope: %w", err)
		}
		if err := json.Unmarshal(raw.Data, out); err != nil {
			return fmt.Errorf("triepayments: decode data: %w", err)
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
	return fmt.Sprintf("triemint: api error %d: %s", e.StatusCode, e.Message)
}

func IsNotFound(err error) bool {
	var apiErr *APIError
	if err == nil {
		return false
	}
	var e *APIError
	if errors.As(err, &e) {
		return e.StatusCode == 404 || errors.Is(e, apiErr)
	}
	return false
}

func IsUnauthorized(err error) bool {
	var e *APIError
	if errors.As(err, &e) {
		return e.StatusCode == 401
	}
	return false
}

func IsConflict(err error) bool {
	var e *APIError
	if errors.As(err, &e) {
		return e.StatusCode == 409
	}
	return false
}
