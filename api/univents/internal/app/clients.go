package app

import (
	"context"
	"time"

	"lib/errx"
	"lib/objectstorage"
	"univents/internal/config"

	idx "sdk/identityx"
)

func SetupIdentityX(cfg config.Config) *idx.Client {
	projectID := cfg.IdxProjectID
	client, err := idx.NewClient(idx.Config{
		BaseURL:   cfg.IdxURL,
		APIKey:    cfg.IdxAPIKey,
		ProjectID: projectID,
		Debug:     cfg.DebugMode,
	})
	if err != nil {
		errx.Exit(err, "error creating identityx client")
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = client.Tokens.GetJWKS(ctx, false)
		if err != nil {
			errx.Exit(err, "error fetching initial JWKS")
		}
	}()
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
