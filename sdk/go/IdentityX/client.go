package idx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"resty.dev/v3"
)

type Config struct {
	BaseURL   string
	APIKey    string
	ProjectID uuid.UUID
	Debug     bool
}

type Client struct {
	http *resty.Client

	projectID uuid.UUID

	Tokens *TokenService
	Actors *ActorService
}

func NewClient(cfg Config) (*Client, error) {
	if cfg.ProjectID == uuid.Nil || cfg.APIKey == "" {
		return nil, &ConfigError{Field: "ProjectID", Message: "required"}
	}

	if cfg.BaseURL == "" {
		return nil, &ConfigError{Field: "BaseURL", Message: "required"}
	}

	// otelhttp.NewTransport injects traceparent from the request context, so
	// calls made under an active span join the caller's distributed trace
	// instead of starting a fresh root trace on the receiving API.
	hc := resty.NewWithClient(&http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}).
		SetBaseURL(strings.TrimRight(cfg.BaseURL, "/")).
		SetTimeout(30*time.Second).
		SetHeader("Accept", "application/json").
		SetHeader("X-API-Key", cfg.APIKey).
		SetDebug(cfg.Debug)

	c := &Client{
		http:      hc,
		projectID: cfg.ProjectID,
	}

	c.Tokens = &TokenService{client: c, cacheTTL: time.Hour}
	c.Actors = &ActorService{client: c}
	return c, nil
}

// Close releases the underlying HTTP client's resources.
func (c *Client) Close() error {
	return c.http.Close()
}

// send executes req and decodes a 2xx body into the request's result. Any
// 4xx/5xx comes back as *APIError carrying the HTTP status, whatever the body
// looks like.
func (c *Client) send(req *resty.Request, method, path string) error {
	res, err := req.SetResponseExpectContentType("application/json").Execute(method, path)
	if err != nil {
		return fmt.Errorf("idx: %s %s: %w", method, path, err)
	}
	if res.IsStatusFailure() {
		apiErr := &APIError{}
		_ = json.Unmarshal(res.Bytes(), apiErr)
		apiErr.StatusCode = res.StatusCode()
		if apiErr.Message == "" {
			apiErr.Message = apiErr.Title
		}
		if apiErr.Message == "" {
			apiErr.Message = res.String()
		}
		return apiErr
	}
	return nil
}

// Bootstrap creates a client and verifies the initial JWKS so the caller can
// serve requests immediately. It blocks up to the given timeout; on failure it
// returns the error so the caller decides how to fail fast.
func Bootstrap(ctx context.Context, cfg Config) (*Client, error) {
	client, err := NewClient(cfg)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err = client.Tokens.GetJWKS(ctx, false)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// MustBootstrap is Bootstrap that panics on failure — for boot-time callers
// that treat an unreachable IdentityX as fatal.
func MustBootstrap(ctx context.Context, cfg Config) *Client {
	client, err := Bootstrap(ctx, cfg)
	if err != nil {
		panic(err)
	}
	return client
}
