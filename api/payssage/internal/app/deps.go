package app

import (
	"context"
	"lib/database"
	"lib/errx"
	"lib/telemetry"
	"lib/validator"
	"net/http"
	"payssage/internal/providers"
	"payssage/ports"
	"time"

	idx "sdk/identityx"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
	fm "github.com/MintzyG/fun/middlewares"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"resty.dev/v3"
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
	database.SetConstraintErrorRegistry(database.ConstraintRegistry{
		// intents
		"chk_intents_amount_cents": "amount must be greater than zero",
		"chk_intents_status":       "invalid intent status",

		// oauth_states
		"chk_oauth_states_flow": "invalid oauth flow type",

		// org_members
		"chk_org_members_role": "invalid organization member role",

		// wallets
		"chk_wallets_fee_bps":        "fee (bps) must be non-negative",
		"uniq_wallets_org_name":      "a wallet with this name already exists in this organization",
		"uniq_wallets_personal_name": "a wallet with this name already exists",

		// webhook_deliveries
		"chk_webhook_deliveries_status": "invalid webhook delivery status",

		// webhook_events
		"uniq_webhook_events_external_id": "a webhook event with this external id already exists for this provider",

		// organizations
		"uniq_organizations_slug": "an organization with this slug already exists",

		// provider_credentials
		"uniq_provider_credentials_active": "credentials for this provider are already connected to this wallet",

		// sellers
		"uniq_sellers_active": "this seller is already connected to this wallet",
	})
}

func SetupIdentityX(cfg Config) *idx.Client {
	client, err := idx.NewClient(idx.Config{
		BaseURL:   cfg.IdxURL,
		APIKey:    cfg.IdxAPIKey,
		ProjectID: cfg.IdxProjectID,
		Debug:     true,
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

func setupProviders(cfg Config) {
	mercadoPago, err := providers.NewMercadoPagoProvider(
		cfg.MpClientID,
		cfg.MpAccessToken,
		cfg.MpClientSecret,
		cfg.MpRedirectURI,
		cfg.MpWebhookSecret,
	)
	if err != nil {
		errx.Exit(err, "error setting up mercado pago")
	}

	providers.PayssageProviders.OAuth = map[providers.AvailableProviders]ports.OAuthProvider{
		providers.MercadoPagoProvider: mercadoPago,
	}

	providers.PayssageProviders.Payments = map[providers.AvailableProviders]ports.PaymentAbstractionLayer{
		providers.MercadoPagoProvider: mercadoPago,
	}
}

func (app *Payssage) setupAuthMiddlewares() *fm.Middleware[*idx.AccessClaims] {
	keyFunc := func(ctx context.Context, tokenStr string) (*idx.AccessClaims, error) {
		return app.idxClient.Tokens.VerifyAccessToken(ctx, tokenStr)
	}

	jwtHook := func(ctx context.Context, claims *idx.AccessClaims) (context.Context, error) {
		return idx.WithIdentity(ctx, &idx.Identity{
			Sub: idx.Subject{
				ID:           claims.Sub.ID,
				ProjectID:    claims.Sub.ProjectID,
				Email:        claims.Sub.Email,
				Type:         claims.Sub.Type,
				Capabilities: claims.Sub.Capabilities,
				Metadata:     claims.Sub.Metadata,
			},
			Cred: idx.Credential{
				Type: "token",
			},
		}), nil
	}

	apiKeyHook := func(ctx context.Context, rawKey string) (context.Context, error) {
		var ident idx.Identity
		_, err := resty.New().R().
			WithContext(ctx).
			SetHeader("X-API-KEY", rawKey).
			SetResult(&ident).
			Get("http://identityx:8080/auth/introspect")
		if err != nil {
			telemetry.Log().Error("error fetching identity", zap.Error(err))
			return ctx, err
		}
		return idx.WithIdentity(ctx, &ident), nil
	}

	return fm.New[*idx.AccessClaims](keyFunc, jwtHook, apiKeyHook)
}
