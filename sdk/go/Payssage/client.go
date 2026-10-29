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
		return &apiErr
	}

	if out != nil {
		err = decode(res.Bytes(), out)
		if err != nil {
			return fmt.Errorf("payssage: decode data: %w", err)
		}
	}

	return nil
}

// decode reads a success body into out, with or without the fun envelope:
// an envelope is a numeric code next to data or module.
func decode(body []byte, out any) error {
	var probe struct {
		Code   *int            `json:"code"`
		Data   json.RawMessage `json:"data"`
		Module *string         `json:"module"`
	}
	if json.Unmarshal(body, &probe) == nil && probe.Code != nil && (probe.Data != nil || probe.Module != nil) {
		if len(probe.Data) == 0 {
			return nil
		}
		body = probe.Data
	}
	return json.Unmarshal(body, out)
}

// list reads a list payload: `{ "items": [...] }`, or the bare array a
// Payssage still on the envelope sends.
type list[T any] []T

func (l *list[T]) UnmarshalJSON(b []byte) error {
	var items []T
	if json.Unmarshal(b, &items) != nil {
		var obj struct {
			Items []T `json:"items"`
		}
		err := json.Unmarshal(b, &obj)
		if err != nil {
			return err
		}
		items = obj.Items
	}
	*l = items
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
