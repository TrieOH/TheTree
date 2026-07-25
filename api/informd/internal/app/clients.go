package app

import (
	"Informd/internal/config"
	"context"
	"lib/errx"
	"time"

	idx "sdk/identityx"
)

func SetupIdentityX(cfg config.Config) *idx.Client {
	client, err := idx.NewClient(idx.Config{
		BaseURL:   cfg.IdxURL,
		APIKey:    cfg.IdxAPIKey,
		ProjectID: cfg.IdxProjectID,
		Debug:     true,
	})
	if err != nil {
		errx.Exit(err, "error creating identity_x client")
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = client.Tokens.GetJWKS(ctx, true)
		if err != nil {
			errx.Exit(err, "error fetching initial JWKS")
		}
	}()
	return client
}
