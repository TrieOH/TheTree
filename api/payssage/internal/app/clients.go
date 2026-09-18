package app

import (
	"context"

	"payssage/internal/config"

	idx "sdk/identityx"

	"time"

	"resty.dev/v3"
)

func SetupIdentityX(ctx context.Context, cfg config.Config) (*idx.Client, error) {
	return idx.Bootstrap(ctx, cfg.ToIdentityXConfig())
}

func SetupHTTPClient() *resty.Client {
	return resty.New().SetTimeout(15 * time.Second)
}
