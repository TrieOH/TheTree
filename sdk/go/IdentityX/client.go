package idx

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/MintzyG/sdkkit"
	"github.com/google/uuid"
)

type Config struct {
	BaseURL            string
	APIKey             string
	ProjectID          uuid.UUID
	Debug              bool
	EncryptionPassword string
	CredsFilePath      string
}

type Client struct {
	*sdkkit.Client
	projectID uuid.UUID
	baseURL   string
	debug     bool

	Creds  *CredentialHandler
	Tokens *TokenService
}

var setupComplete atomic.Bool

func NewClient(cfg Config) (*Client, error) {
	if cfg.ProjectID == uuid.Nil && cfg.APIKey == "" {
		// Allow zero-config construction - caller will finish via Setup().
	} else if cfg.ProjectID == uuid.Nil {
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

	if cfg.ProjectID != uuid.Nil && cfg.APIKey != "" {
		setupComplete.Store(true)
	} else {
		setupComplete.Store(false)
	}

	c.Creds = NewCredentialHandler(cfg.CredsFilePath, []byte(cfg.EncryptionPassword))
	c.Tokens = &TokenService{client: c, cacheTTL: time.Hour}
	return c, nil
}

// IsSetupComplete reports whether the client has been fully configured with
// an API key and project ID — either at construction or via Setup.
func (c *Client) IsSetupComplete() bool {
	return setupComplete.Load()
}

// Setup finishes client configuration after the IdentityX project-setup flow
// returns an API key and project ID. It persists the credentials via the
// attached CredentialHandler (if any), reconstructs the inner HTTP client with
// the API key, and marks the client as ready.
//
//	handler := idx.NewCredentialHandler("./creds.enc", password)
//	client.Creds = handler
//	client.Setup(apiKey, projectID)
func (c *Client) Setup(apiKey string, projectID uuid.UUID) error {
	if apiKey == "" || projectID == uuid.Nil {
		return &ConfigError{Field: "Setup", Message: "apiKey and projectID are required"}
	}

	if c.Creds != nil {
		if err := c.Creds.SaveCreds(apiKey, projectID); err != nil {
			return fmt.Errorf("idx: setup: save creds: %w", err)
		}
	}

	core, err := sdkkit.New(sdkkit.Config{
		BaseURL: c.baseURL,
		APIKey:  apiKey,
		Debug:   c.debug,
	})
	if err != nil {
		return fmt.Errorf("idx: setup: %w", err)
	}

	c.Client = core
	c.projectID = projectID
	setupComplete.Store(true)
	return nil
}
