package app

import (
	"context"

	"payssage/internal/config"

	idx "sdk/identityx"

	"time"

	"resty.dev/v3"
)

func SetupIdentityX(cfg config.Config) *idx.Client {
	return idx.MustBootstrap(context.Background(), cfg.ToIdentityXConfig())
}

func SetupHTTPClient() *resty.Client {
	return resty.New().SetTimeout(15 * time.Second)
}
