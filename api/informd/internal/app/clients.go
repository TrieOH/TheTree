package app

import (
	"context"

	"Informd/internal/config"

	idx "sdk/identityx"
)

func SetupIdentityX(cfg config.Config) *idx.Client {
	return idx.MustBootstrap(context.Background(), cfg.ToIdentityXConfig())
}
