package initialization

import (
	http "TriePayments/internal/shared/validation"
	"log"
	"strings"

	"github.com/spf13/viper"
)

func LoadEnv(app *TriePayments) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := http.LoadProxyConfig(); err != nil {
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

	Port := viper.GetString("PORT")
	if Port == "" {
		Port = "8080"
	}
	app.Port = Port
}
