package app

import (
	"context"

	"lib/errx"
	"lib/objectstorage"
	"univents/internal/config"

	idx "sdk/identityx"
)

func SetupIdentityX(cfg config.Config) *idx.Client {
	client, err := idx.Bootstrap(context.Background(), idx.Config{
		BaseURL:   cfg.IdxURL,
		APIKey:    cfg.IdxAPIKey,
		ProjectID: cfg.IdxProjectID,
		Debug:     cfg.DebugMode,
	})
	if err != nil {
		errx.Exit(err, "error creating identityx client")
	}
	return client
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
