package app

import (
	"Informd/internal/config"
	"context"
	"lib/errx"

	idx "sdk/identityx"
)

func SetupIdentityX(cfg config.Config) *idx.Client {
	client, err := idx.Bootstrap(context.Background(), idx.Config{
		BaseURL:   cfg.IdxURL,
		APIKey:    cfg.IdxAPIKey,
		ProjectID: cfg.IdxProjectID,
		Debug:     true,
	})
	if err != nil {
		errx.Exit(err, "error creating identityx client")
	}
	return client
}
