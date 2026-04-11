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
	if viper.GetString("GOAUTH_URL") == "" {
		log.Fatal("GOAUTH_URL must be set")
	}
	if viper.GetString("GOAUTH_API_KEY") == "" {
		log.Fatal("GOAUTH_API_KEY must be set")
	}
	if viper.GetString("GO_AUTH_PROJECT_ID") == "" {
		log.Fatal("GO_AUTH_PROJECT_ID must be set")
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
