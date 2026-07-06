package idx

import (
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
