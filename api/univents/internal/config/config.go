package config

import (
	"uuid"

	libconfig "lib/config"
)

type Config struct {
	libconfig.Server
	libconfig.IdentityX
	libconfig.CORS

	Postgres libconfig.Postgres `envPrefix:"UNIVENTS_"`
	SMTP     libconfig.SMTP
	Jobs     libconfig.Jobs

	// Feature fields
	AppURL string `env:"APP_URL,required"`

	// When true, every JWT-authenticated route requires the actor's email to
	// be verified (identityx subject.verified_at set). Unverified actors get
	// 403. Public routes (webhooks, health, ws) are unaffected.
	RequireVerifiedEmail bool `env:"REQUIRE_VERIFIED_EMAIL" envDefault:"false"`

	// Payssage
	PayssageURL           string    `env:"PAYSSAGE_URL,required"`
	PayssageAPIKey        string    `env:"PAYSSAGE_API_KEY,required"`
	PayssageWalletID      uuid.UUID `env:"PAYSSAGE_WALLET_ID,required"` // one platform wallet shared by every event (D6)
	PayssageWebhookSecret string    `env:"PAYSSAGE_WEBHOOK_SECRET,required"`
	// Off on Cloudflare: the check wakes Payssage (and IdentityX) on every
	// cold start.
	PayssageVerifyWalletOnBoot bool `env:"PAYSSAGE_VERIFY_WALLET_ON_BOOT" envDefault:"true"`

	// Object Storage (RustFS)
	ObjStorageEndpoint  string `env:"OBJECT_STORAGE_ENDPOINT,required"`
	ObjStorageAccessKey string `env:"OBJECT_STORAGE_ACCESS_KEY,required"`
	ObjStorageSecretKey string `env:"OBJECT_STORAGE_SECRET_KEY,required"`
	ObjStorageUseSSL    bool   `env:"OBJECT_STORAGE_USE_SSL"             envDefault:"true"`
	ObjStorageRegion    string `env:"OBJECT_STORAGE_REGION"              envDefault:"us-east-1"`

	// Webhook HMAC
	HmacSecret string `env:"HMAC_SECRET,required"`
}
