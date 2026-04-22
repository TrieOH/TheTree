package app

import (
	"TrieForms/internal/shared/validation"
	"log"
	"strings"

	"github.com/spf13/viper"
)

func LoadEnv() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := validation.LoadProxyConfig(); err != nil {
		log.Fatalf("LoadProxyConfig failed: %v", err.Error())
	}
	if viper.GetString("IDENTITY_X_URL") == "" {
		log.Fatal("IDENTITY_X_URL must be set")
	}
	if viper.GetString("IDENTITY_X_API_KEY") == "" {
		log.Fatal("IDENTITY_X_API_KEY must be set")
	}
	if viper.GetString("IDENTITY_X_PROJECT_ID") == "" {
		log.Fatal("IDENTITY_X_PROJECT_ID must be set")
	}
	if viper.GetString("SPICEDB_ADDR") == "" {
		log.Fatal("SPICEDB_ADDR must be set")
	}
	if viper.GetString("SPICEDB_TOKEN") == "" {
		log.Fatal("SPICEDB_TOKEN must be set")
	}
	if viper.GetString("PORT") == "" {
		log.Fatal("PORT must be set")
	}
}
