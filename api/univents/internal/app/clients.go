package app

import (
	"context"

	"lib/errx"
	"lib/objectstorage"
	"univents/internal/config"

	idx "sdk/identityx"
	payssage "sdk/payssage"

	"github.com/google/uuid"
)

func SetupIdentityX(cfg config.Config) *idx.Client {
	return idx.MustBootstrap(context.Background(), cfg.ToIdentityXConfig())
}

// SetupPayssage builds the service-to-service Payssage client. Calls are
// authenticated with the platform API key, so wallets/sellers are owned by
// the owner of Univents (the platform identity). The provider OAuth callback
// URL is Payssage's own concern (D7) — univents' `PAYSSAGE_URL` is
// server-to-server only.
func SetupPayssage(cfg config.Config) *payssage.Client {
	return payssage.New(payssage.Config{
		BaseURL: cfg.PayssageURL,
		APIKey:  cfg.PayssageAPIKey,
	})
}

// VerifyPayssageWallet is the fail-fast boot check for the platform wallet
// (D6): every event's seller lives on this one precreated wallet, so a wrong
// or unreachable wallet id means the store cannot work — do not start.
func VerifyPayssageWallet(ctx context.Context, client *payssage.Client, walletID uuid.UUID) {
	_, err := client.GetWallet(ctx, walletID)
	if err != nil {
		errx.Exit(err, "payssage platform wallet check failed (PAYSSAGE_WALLET_ID)")
	}
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
