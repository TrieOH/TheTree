package app

import (
	"context"
	"lib/errx"
	"lib/validator"
	"log"
	"net/http"
	"payssage/internal/platform/providers"
	"payssage/ports"
	"time"

	idx "sdk/identityx"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
	"github.com/go-chi/chi/v5"
)

func SetupFUN(module string) {
	fun.SetConfig(fun.Config{
		MaxTraceSize:         50,
		ResponseSizeLimit:    10 * 1024 * 1024,
		MaxInterceptorAmount: 20,
		DefaultContentType:   "application/json",
		EnableSizeValidation: true,
		DefaultModule:        module,
	})

	v := validator.SetupValidator()
	bind.SetValidator(v)
	fun.SetPathParamFunc(func(r *http.Request, key string) string {
		return chi.URLParam(r, key)
	})
}

func SetupConstraintMessages() {
	//database.SetConstraintErrorRegistry(database.ConstraintRegistry{})
}

func SetupIdentityX(cfg Config) *idx.Client {
	client, err := idx.NewClient(idx.Config{
		BaseURL:   cfg.IdxURL,
		APIKey:    cfg.IdxAPIKey,
		ProjectID: cfg.IdxProjectID,
		Debug:     cfg.DebugMode,
	})
	if err != nil {
		errx.Exit(err, "error creating identity_x client")
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = client.Tokens.GetJWKS(ctx, true)
		if err != nil {
			errx.Exit(err, "error fetching initial JWKS")
		}
	}()
	return client
}

func setupPaymentProviders(cfg Config) paymentProviders {
	var pp paymentProviders
	mpProvider, err := providers.NewMercadoPagoProvider(
		cfg.MpClientID,
		cfg.MpAccessToken,
		cfg.MpClientSecret,
		cfg.MpRedirectURI,
		cfg.MpWebhookSecret,
	)
	if err != nil {
		log.Fatalf("Error creating mercado pago provider: %s", err.Error())
	}

	pp.oauth = map[string]ports.OAuthProvider{
		"mercadopago": mpProvider,
	}

	pp.payments = map[string]ports.PaymentAbstractionLayer{
		"mercadopago": mpProvider,
	}

	return pp
}
