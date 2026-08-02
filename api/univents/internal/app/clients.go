package app

import (
	"context"

	"lib/errx"
	"lib/objectstorage"
	"univents/internal/config"

	idx "sdk/identityx"
)

func SetupIdentityX(cfg config.Config) *idx.Client {
	return idx.MustBootstrap(context.Background(), cfg.ToIdentityXConfig())
}

func SetupObjectStorage(cfg config.Config) *objectstorage.Client {
	client, err := objectstorage.New(context.Background(), objectstorage.Config{
		Endpoint:  cfg.ObjStorageEndpoint,
		AccessKey: cfg.ObjStorageAccessKey,
		SecretKey: cfg.ObjStorageSecretKey,
		UseSSL:    cfg.ObjStorageUseSSL,
		Region:    cfg.ObjStorageRegion,
	})
	if err != nil {
		errx.Exit(err, "failed to create object-storage client")
	}
	return client
}
