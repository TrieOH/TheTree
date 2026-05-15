package app

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

func LoadEnv() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if viper.GetString("IDENTITY_X_URL") == "" {
		log.Fatal("IDENTITY_X_URL must be set")
	}
	if viper.GetString("IDENTITY_X_API_KEY") == "" {
		log.Fatal("IDENTITY_X_API_KEY must be set")
	}
	if viper.GetString("IDENTITY_X_PROJECT_ID") == "" {
		log.Fatal("IDENTITY_X_PROJECT_ID must be set")
	}

	if viper.GetString("MP_CLIENT_ID") == "" {
		log.Fatal("MP_CLIENT_ID must be set")
	}
	if viper.GetString("MP_CLIENT_SECRET") == "" {
		log.Fatal("MP_CLIENT_SECRET must be set")
	}
	if viper.GetString("MP_ACCESS_TOKEN") == "" {
		log.Fatal("MP_ACCESS_TOKEN must be set")
	}
	if viper.GetString("MP_REDIRECT_URI") == "" {
		log.Fatal("MP_REDIRECT_URI must be set")
	}
	if viper.GetString("MP_WEBHOOK_SECRET") == "" {
		log.Fatal("MP_WEBHOOK_SECRET must be set")
	}
	if viper.GetString("MP_TEST_ACCESS_TOKEN") == "" {
		log.Fatal("MP_TEST_ACCESS_TOKEN must be set")
	}
	if viper.GetString("MP_TEST_PUBLIC_KEY") == "" {
		log.Fatal("MP_TEST_PUBLIC_KEY must be set")
	}

	if viper.GetString("PORT") == "" {
		log.Fatal("PORT must be set")
	}
}
