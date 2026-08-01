package app

import (
	"context"
	"lib/errx"
	"payssage/internal/config"
	"time"

	idx "sdk/identityx"

	"resty.dev/v3"
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

func SetupHTTPClient() *resty.Client {
	return resty.New().SetTimeout(15 * time.Second)
}
