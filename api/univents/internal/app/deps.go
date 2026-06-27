package app

import (
	"context"
	"lib/validator"
	"log"
	"net/http"
	"time"

	"lib/errx"
	"lib/telemetry"

	idx "sdk/identityx"
	"sdk/payssage"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
	mws "github.com/MintzyG/fun/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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
	//database.SetConstraintErrorRegistry(database.ConstraintRegistry{})
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

func SetupPayssage(cfg Config) *payssage.Client {
	paymentsURL := cfg.PayssageURL
	paymentsAPIKey := cfg.PayssageAPIKey
	client := payssage.New(paymentsURL, paymentsAPIKey)
	return client
}

func SetupObjectStorage(cfg Config) *minio.Client {
	client, err := minio.New(cfg.ObjStorageURL, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.ObjStorageAccessKey, cfg.ObjStorageSecretKey, ""),
		Secure: cfg.ObjStorageUseSSL,
	})
	if err != nil {
		errx.Exit(err, "failed to create ObjectStorage client")
	}
	log.Println("Connected to ObjectStorage")
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
