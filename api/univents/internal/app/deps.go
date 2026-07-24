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
		"chk_event_status_valid":             "Event status must be one of: draft, active, discontinued.",
		"chk_event_payments_config_complete": "Both Payssage seller ID and wallet ID must be set together, or neither.",

		// event_members
		"chk_event_members_role_valid": "Member role must be one of: owner, admin, staff.",

		// editions
		"chk_editions_dates_valid":               "Edition end date must be after the start date.",
		"chk_editions_registration_before_start": "Registration opening date must be before or equal to the edition start date.",
		"excl_editions_no_overlap":               "This edition's dates overlap with another edition of the same event.",

		// registrations
		"chk_registrations_status_valid": "Registration status must be one of: pending, confirmed, cancelled, expired.",

		// products
		"uniq_products_edition_vendor_code":         "A product with this vendor code already exists in this edition.",
		"uniq_product_variants_edition_vendor_code": "A product variant with this vendor code already exists in this edition.",

		// product_purchases
		"chk_product_purchases_status_valid": "Product purchase status must be one of: pending, confirmed, cancelled, expired.",

		// programs
		"chk_programs_kind_valid": "Program kind must be one of: activity, checkpoint.",

		// program_occurrences
		"chk_program_occurrences_dates_valid": "Program occurrence end time must be after the start time.",

		// program_participations
		"chk_program_participations_status_valid": "Participation status must be one of: registered, attended, no_show, cancelled.",

		// signatures
		"chk_signatures_status_valid":    "Signature status must be one of: requested, ready, declined, expired.",
		"chk_signatures_ready_has_image": "A signature marked as ready must have an image URL set.",

		// certification_templates
		"chk_certification_templates_kind_valid": "Certification template kind must be one of: edition_attendance, program_attendance.",
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
