package idx

import (
	"context"
	"time"

	"github.com/MintzyG/sdkkit"
	"github.com/google/uuid"
)

type Config struct {
	BaseURL   string
	APIKey    string
	ProjectID uuid.UUID
	Debug     bool
}

type Client struct {
	*sdkkit.Client

	projectID uuid.UUID
	baseURL   string
	debug     bool

	Tokens *TokenService
	Actors *ActorService
}

func NewClient(cfg Config) (*Client, error) {
	if cfg.ProjectID == uuid.Nil || cfg.APIKey == "" {
		return nil, &ConfigError{Field: "ProjectID", Message: "required"}
	}

	core, err := sdkkit.New(sdkkit.Config{
		BaseURL: cfg.BaseURL,
		APIKey:  cfg.APIKey,
		Debug:   cfg.Debug,
	})
	if err != nil {
		return nil, err
	}

	c := &Client{
		Client:    core,
		projectID: cfg.ProjectID,
		baseURL:   cfg.BaseURL,
		debug:     cfg.Debug,
	}

	c.Tokens = &TokenService{client: c, cacheTTL: time.Hour}
	c.Actors = &ActorService{client: c}
	return c, nil
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

	if _, err := client.Tokens.GetJWKS(ctx, false); err != nil {
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
