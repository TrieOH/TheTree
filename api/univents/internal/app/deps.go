package app

import (
	"context"
	"lib/database"
	"lib/errx"
	"lib/objectstorage"
	"lib/telemetry"
	"lib/validator"
	"net/http"
	"time"

	idx "sdk/identityx"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
	mws "github.com/MintzyG/fun/middlewares"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type SimpleLogger struct{}

func (l *SimpleLogger) Intercept(_ context.Context, rs *fun.Response, statusCode int) {
	if statusCode == 500 {
		telemetry.Log().Info("InternalServerError Response", zap.Any("response", rs))
	}
}

func (l *SimpleLogger) InterceptSimple(rs *fun.Response, statusCode int) {
	if statusCode == 500 {
		telemetry.Log().Info("InternalServerError Response", zap.Any("response", rs))
	}
}

func SetupFUN(module string) {
	fun.SetConfig(fun.Config{
		MaxTraceSize:         50,
		ResponseSizeLimit:    10 * 1024 * 1024,
		MaxInterceptorAmount: 20,
		DefaultContentType:   "application/json",
		EnableSizeValidation: true,
		DefaultModule:        module,
	})

	err := fun.AddInterceptor(&SimpleLogger{})
	if err != nil {
		errx.Exit(err, "failed to add interceptor")
	}

	v := validator.SetupValidator()
	bind.SetValidator(v)
	fun.SetPathParamFunc(func(r *http.Request, key string) string {
		return chi.URLParam(r, key)
	})
}

func SetupConstraintMessages() {
	database.SetConstraintErrorRegistry(database.ConstraintRegistry{
		// events
		"chk_series_requires_single_edition": "a non-series event can only have a single edition",

		// editions
		"editions_check":  "registration close date must be after the registration open date",
		"editions_check1": "edition end date must be after the start date",

		// tickets
		"tickets_name_edition_id_key": "a ticket with this name already exists for this edition",

		// activities
		"activities_check":  "activity end date must be after the start date",
		"activities_check1": "capacity must be greater than zero when capacity is enabled",
		"activities_check2": "capacity cannot be lower than the remaining capacity",
		"activity_interest_list_activity_id_user_id_key": "user is already on the interest list for this activity",

		// products
		"chk_product_status": "invalid product status",
		"products_check":     "a ticket-type product must reference a ticket",
		"products_check1":    "inventory quantity must be greater than zero when inventory is enabled",
		"products_check2":    "inventory quantity cannot be lower than the remaining inventory",
		"product_bundle_components_component_type_check":                  "component type must be either 'ticket' or 'product'",
		"product_bundle_components_bundle_product_id_component_type__key": "this component is already part of the bundle",
		"product_reservations_session_id_product_id_key":                  "a reservation for this product already exists for this session",

		// purchases
		"uq_purchases_session_id": "a purchase already exists for this session",

		// ticket permissions
		"chk_type_matches_target": "the permission target does not match the selected permission type",

		// tokens
		"user_token_balances_user_id_edition_id_key": "a token balance already exists for this user in this edition",

		// edition registrations / interest
		"edition_interest_list_edition_id_user_id_key": "user is already on the interest list for this edition",
		"edition_registrations_edition_id_user_id_key": "user is already registered for this edition",
	})
}

func SetupIdentityX(cfg Config) *idx.Client {
	projectID := cfg.IdxProjectID
	client, err := idx.NewClient(idx.Config{
		BaseURL:   cfg.IdxURL,
		APIKey:    cfg.IdxAPIKey,
		ProjectID: projectID,
		Debug:     cfg.DebugMode,
	})
	if err != nil {
		errx.Exit(err, "error creating identityx client")
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = client.Tokens.GetJWKS(ctx, false)
		if err != nil {
			errx.Exit(err, "error fetching initial JWKS")
		}
	}()
	return client
}

//func SetupPayssage(cfg Config) *payssage.Client {
//	paymentsURL := cfg.PayssageURL
//	paymentsAPIKey := cfg.PayssageAPIKey
//	client := payssage.New(paymentsURL, paymentsAPIKey)
//	return client
//}

func SetupObjectStorage(cfg Config) *objectstorage.Client {
	client, err := objectstorage.New(context.Background(), objectstorage.Config{
		Endpoint:  cfg.ObjStorageEndpoint,
		AccessKey: cfg.ObjStorageAccessKey,
		SecretKey: cfg.ObjStorageSecretKey,
		UseSSL:    cfg.ObjStorageUseSSL,
		Region:    cfg.ObjStorageRegion,
	})
	if err != nil {
		errx.Exit(err, "failed to create object-storage client")
	}
	return client
}

func SetupAuthMiddlewares() *mws.Middleware[*idx.AccessClaims] {
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
		telemetry.Log().Info("user tried to use api key",
			zap.String("message", "this service does not provide a public api"),
			zap.String("key", rawKey),
		)
		return ctx, fun.ErrForbidden("this service does not provide public access to the api")
	}

	return mws.New[*idx.AccessClaims](keyFunc, jwtHook, apiKeyHook)
}
