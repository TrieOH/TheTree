package payssage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"resty.dev/v3"
)

type Config struct {
	BaseURL string // Payssage API base URL
	APIKey  string // IdentityX access API key resolving to the platform/service actor
}

type Client struct {
	http *resty.Client
}

func New(cfg Config) *Client {
	// otelhttp.NewTransport injects traceparent from the request context, so
	// calls made under an active span join the caller's distributed trace.
	return &Client{
		http: resty.NewWithClient(&http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}).
			SetBaseURL(strings.TrimRight(cfg.BaseURL, "/")).
			SetTimeout(15*time.Second).
			SetHeader("Accept", "application/json").
			SetHeader("X-API-Key", cfg.APIKey),
	}
}

// Close releases the underlying HTTP client's resources.
func (c *Client) Close() error {
	return c.http.Close()
}

func (c *Client) do(ctx context.Context, method, path string, body, out any) error {
	req := c.http.R().SetContext(ctx).SetResponseExpectContentType("application/json")
	if body != nil {
		req.SetBody(body)
	}
	res, err := req.Execute(method, path)
	if err != nil {
		return fmt.Errorf("payssage: %s %s: %w", method, path, err)
	}

	if res.IsStatusFailure() {
		var apiErr APIError
		_ = json.Unmarshal(res.Bytes(), &apiErr)
		apiErr.StatusCode = res.StatusCode()
		if apiErr.Message == "" {
			apiErr.Message = apiErr.Title
		}
		return &apiErr
	}

	if out != nil {
		err = json.Unmarshal(res.Bytes(), out)
		if err != nil {
			return fmt.Errorf("payssage: decode: %w", err)
		}
	}

	return nil
}

// list reads a list payload: `{ "items": [...] }`.
type list[T any] struct {
	Items []T `json:"items"`
}

// APIError is returned when Payssage responds with a failure status.
// StatusCode is the HTTP status; the rest is the decoded RFC 9457 problem,
// which may be empty when the body was not a problem. Message is the
// problem's detail.
type APIError struct {
	StatusCode int    `json:"-"`
	Type       string `json:"type"`
	Title      string `json:"title"`
	Message    string `json:"detail"`
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
