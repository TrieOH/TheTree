package config

import (
	libconfig "lib/config"
)

type Config struct {
	MercadoPagoConfig

	libconfig.Server
	libconfig.IdentityX
	libconfig.CORS

	Postgres libconfig.Postgres `envPrefix:"PAYSSAGE_"`
	Jobs     libconfig.Jobs

	// TestMode turns on the /testmode routes, which fake provider outcomes.
	// Never on in production.
	TestMode bool `env:"TEST_MODE"`
}

type MercadoPagoConfig struct {
	MpClientID        string `env:"MP_CLIENT_ID,required"`
	MpClientSecret    string `env:"MP_CLIENT_SECRET,required"`
	MpAccessToken     string `env:"MP_ACCESS_TOKEN,required"`
	MpRedirectURI     string `env:"MP_REDIRECT_URI,required"`
	MpWebhookSecret   string `env:"MP_WEBHOOK_SECRET,required"`
	MpTestAccessToken string `env:"MP_TEST_ACCESS_TOKEN,required"`
	MpTestPublicKey   string `env:"MP_TEST_PUBLIC_KEY,required"`
}
