package idx

import (
	"time"

	"github.com/TrieOH/sdkkit"
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

	Users  *UserService
	Tokens *TokenService
}

func NewClient(cfg Config) (*Client, error) {
	if cfg.ProjectID == uuid.Nil {
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
	}
	c.Users = &UserService{client: c}
	c.Tokens = &TokenService{client: c, cacheTTL: time.Hour}
	return c, nil
}
