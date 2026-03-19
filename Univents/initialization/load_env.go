package initialization

import (
	"log"
	"strings"
	http2 "univents/internal/shared/validation"

	"github.com/spf13/viper"
)

func LoadEnv(app *UniventsApp) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := http2.LoadProxyConfig(); err != nil {
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

	if viper.GetString("WS_JWT_SECRET") == "" {
		log.Fatal("WS_JWT_SECRET must be set")
	}

	Port := viper.GetString("PORT")
	if Port == "" {
		Port = "8080"
	}
	app.Port = Port
}
