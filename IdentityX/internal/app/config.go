package app

import (
	http2 "IdentityX/internal/shared/validation"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func LoadEnv() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := http2.LoadProxyConfig(); err != nil {
		log.Fatalf("LoadProxyConfig failed: %v", err.Error())
	}
	if iss := viper.GetString("ISSUER"); iss == "" {
		log.Fatalf("ISSUER environment variable not set.")
	}
	if smtpHost := viper.GetString("SMTP_HOST"); smtpHost == "" {
		log.Fatalf("SMTP_HOST environment variable not set.")
	}
	if smtpPort := viper.GetString("SMTP_PORT"); smtpPort == "" {
		log.Fatalf("SMTP_PORT environment variable not set.")
	}

	if keyDuration := viper.GetDuration("IDENTITY_X_KEY_LIFETIME"); keyDuration == time.Duration(0) {
		log.Fatalf("IDENTITY_X_KEY_LIFETIME environment variable not set.")
	}
	if keyDuration := viper.GetDuration("ROTATE_KEYS_JOB_DURATION"); keyDuration == time.Duration(0) {
		log.Fatalf("ROTATE_KEYS_JOB_DURATION environment variable not set.")
	}

	env := viper.GetString("ENV")
	if env == "production" {
		if smtpUsername := viper.GetString("SMTP_USERNAME"); smtpUsername == "" {
			log.Fatalf("SMTP_USERNAME environment variable not set.")
		}
		if smtpPassword := viper.GetString("SMTP_PASSWORD"); smtpPassword == "" {
			log.Fatalf("SMTP_PASSWORD environment variable not set.")
		}
	}
	if smtpFrom := viper.GetString("SMTP_FROM"); smtpFrom == "" {
		log.Fatalf("SMTP_FROM environment variable not set.")
	}
	if viper.GetString("PORT") == "" {
		log.Fatalf("PORT environment variable not set.")
	}
}
