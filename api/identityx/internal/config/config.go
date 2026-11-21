package config

import (
	"time"

	libconfig "lib/config"
)

type Config struct {
	libconfig.Server
	libconfig.CORS

	Postgres libconfig.Postgres `envPrefix:"IDX_"`
	SMTP     libconfig.SMTP
	Jobs     libconfig.Jobs

	// Server extras
	AppURL string `env:"APP_URL,required"`

	// Auth / crypto
	Issuer                string        `env:"ISSUER,required"`
	EncryptionKey         string        `env:"ENCRYPTION_KEY,required"`
	HmacSecret            string        `env:"HMAC_SECRET,required"`
	KeyLifetime           time.Duration `env:"IDENTITY_X_KEY_LIFETIME,required"`
	RotateKeysJobDuration time.Duration `env:"ROTATE_KEYS_JOB_DURATION,required"`

	// Platform-level OAuth logins (IdentityX's own login page). A provider
	// with an empty client id or secret is simply not offered.
	GoogleClientID     string `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
	GoogleRedirectURI  string `env:"GOOGLE_REDIRECT_URI"`
	GithubClientID     string `env:"GITHUB_CLIENT_ID"`
	GithubClientSecret string `env:"GITHUB_CLIENT_SECRET"`
	GithubRedirectURI  string `env:"GITHUB_REDIRECT_URI"`

	// Tokens
	AccessTokenLifetime  time.Duration `env:"ACCESS_TOKEN_LIFETIME,required"`
	RefreshTokenLifetime time.Duration `env:"REFRESH_TOKEN_LIFETIME,required"`

	// Action tokens (email verify / password reset links)
	EmailVerifyTokenTTL time.Duration `env:"EMAIL_VERIFY_TOKEN_TTL" envDefault:"10m"`
	EmailResetTokenTTL  time.Duration `env:"EMAIL_RESET_TOKEN_TTL"  envDefault:"10m"`
}
