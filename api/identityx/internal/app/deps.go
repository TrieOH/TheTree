package app

import (
	"IdentityX/internal/database/sqlc"
	"IdentityX/internal/shared/authz"
	"IdentityX/internal/shared/ports"
	"IdentityX/internal/shared/security"
	"IdentityX/models"
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"lib/crypto"
	"lib/database"
	"lib/errx"
	"lib/telemetry"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"
	"unicode"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
	"github.com/MintzyG/fun/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-co-op/gocron/v2"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func SetupValidator() *validator.Validate {
	var v = validator.New()
	if err := v.RegisterValidation("uuid7", func(fl validator.FieldLevel) bool {
		vv := fl.Field().String()

		u, err := uuid.Parse(vv)
		if err != nil {
			return false
		}

		return u.Version() == 7
	}); err != nil {
		errx.Exit(err, "failed to register uuid7 validator")
	}

	// Custom password validation - requires uppercase, number, and symbol
	if err := v.RegisterValidation("passwd", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		var hasUpper, hasNumber, hasSymbol bool

		for _, c := range password {
			switch {
			case unicode.IsUpper(c):
				hasUpper = true
			case unicode.IsNumber(c):
				hasNumber = true
			case unicode.IsPunct(c) || unicode.IsSymbol(c):
				hasSymbol = true
			}
		}

		return hasUpper && hasNumber && hasSymbol
	}); err != nil {
		errx.Exit(err, "failed to register passwd validator")
	}

	// Use JSON field names for better API responses
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		if name == "" {
			return fld.Name
		}
		return name
	})

	return v
}

func SetupFUN() {
	fun.SetConfig(fun.Config{
		MaxTraceSize:         50,
		ResponseSizeLimit:    10 * 1024 * 1024,
		MaxInterceptorAmount: 20,
		DefaultContentType:   "application/json",
		EnableSizeValidation: true,
		DefaultModule:        "IdentityX-API",
	})

	v := SetupValidator()
	bind.SetValidator(v)
	fun.SetPathParamFunc(func(r *http.Request, key string) string {
		return chi.URLParam(r, key)
	})
}

func SetupConstraintMessages() {
	database.SetConstraintErrorRegistry(database.ConstraintRegistry{
		"chk_valid_user_type":   "user type must be one of: client, project",
		"one_email_for_client":  "an account with this email already exists",
		"one_email_per_project": "an account with this email already exists in this project",

		// sessions
		"chk_session_valid_user_type":           "session user type must be one of: client, project",
		"chk_session_not_revoked_before_issued": "a session cannot be revoked before it was issued",
		"sessions_token_id_key":                 "a session with this token ID already exists",

		// key_pair
		"chk_key_pair_key_type_valid":                 "key type must be one of: goauth, project",
		"chk_key_pair_usage_valid":                    "key usage must be one of: sign, verify",
		"chk_key_pair_status_valid":                   "key status must be one of: active, rotated, revoked",
		"chk_key_pair_type_project_consistency_check": "goauth keys must not have a project, project keys must have a project",
		"chk_key_pair_cant_sign_if_rotated":           "a rotated key pair cannot be used for signing",
		"key_pair_kid_key":                            "a key pair with this kid already exists",
		"one_identity_x_active_signing_key":           "there can only be one active goauth signing key",
	})
}

func SetupRuntimeEnv(db *pgxpool.Pool, encryptionKey []byte) {
	q := sqlc.New(db)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// First, rotate any expired keys to clear the way for new key creation
	if err := q.RotateExpiredKeys(ctx); err != nil {
		log.Printf("Warning: failed to rotate expired signing keys: %v", err)
	}

	// Also run the full key rotation logic to create new keys for projects without active keys
	if err := tryRotateIDXKeys(ctx, q, encryptionKey); err != nil {
		log.Printf("Warning: failed to rotate goauth keys: %v", err)
	}

	if err := tryRotateProjectKeys(ctx, q, encryptionKey); err != nil {
		log.Printf("Warning: failed to rotate project keys: %v", err)
	}

	_, err := q.GetActiveSigningKey(ctx, nil)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			var pub string
			var priv []byte
			pub, priv, err = crypto.GenerateEd25519()
			if err != nil {
				log.Fatalf("failed to generate GoAuth key: %v", err)
			}
			defer zero(priv)

			var encryptedPriv []byte
			encryptedPriv, err = crypto.Encrypt(priv, encryptionKey)
			if err != nil {
				log.Fatalf("failed to encrypt GoAuth key: %v", err)
			}

			keyLifetime := errx.MustEnv("IDENTITY_X_KEY_LIFETIME", time.ParseDuration)

			kid := "goauth:" + ulid.Make().String()
			expiresAt := time.Now().Add(keyLifetime)

			_, err = q.CreateKeyPair(ctx, sqlc.CreateKeyPairParams{
				Kid:             kid,
				ProjectID:       nil,
				KeyType:         "goauth",
				Algorithm:       "EdDSA",
				PublicKey:       pub,
				PrivateKey:      encryptedPriv,
				Usage:           "sign",
				Status:          "active",
				ExpiresAt:       expiresAt,
				VerifyExpiresAt: expiresAt.Add(keyLifetime),
			})

			if err != nil {
				if fun.Is(err, fun.CodeConflict) {
					log.Println("GoAuth signing key already created by another instance")
				} else {
					log.Fatalf("failed to create GoAuth signing key: %s", err)
				}
			} else {
				log.Println("Created GoAuth signing key")
			}
		} else {
			log.Fatalf("failed checking GoAuth signing key: %v", err.Error())
		}
	}

	if err != nil {
		if fun.Is(err, fun.CodeConflict) {
			log.Println("GoAuth Global scope already created by another instance")
		} else {
			log.Fatalf("Failed to create GoAuth Global scope: %s", err)
		}
	} else {
		log.Println("Created GoAuth Global scope")
	}
}

func tryRotateIDXKeys(ctx context.Context, q *sqlc.Queries, encryptionKey []byte) error {
	dbe := database.NewErrorHandler("idx key")
	key, err := q.GetActiveSigningKey(ctx, nil)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			// defensive: no signing key → create
			return createIDXKey(ctx, q, encryptionKey)
		}
		return dbe(err)
	}

	if time.Until(key.ExpiresAt) > 24*time.Hour {
		return nil
	}

	if err = q.RotateSigningKeys(ctx, nil); err != nil {
		return dbe(err)
	}

	return createIDXKey(ctx, q, encryptionKey)
}

func createIDXKey(ctx context.Context, q *sqlc.Queries, encryptionKey []byte) error {
	pub, priv, err := crypto.GenerateEd25519()
	defer zero(priv)
	if err != nil {
		return err
	}

	encryptedPriv, err := crypto.Encrypt(priv, encryptionKey)
	if err != nil {
		return err
	}

	keyLifetime := errx.MustEnv("IDENTITY_X_KEY_LIFETIME", time.ParseDuration)
	kid := "goauth:" + ulid.Make().String()
	expiresAt := time.Now().Add(keyLifetime)

	_, err = q.CreateKeyPair(ctx, sqlc.CreateKeyPairParams{
		Kid:             kid,
		ProjectID:       nil,
		KeyType:         "goauth",
		Algorithm:       "EdDSA",
		PublicKey:       pub,
		PrivateKey:      encryptedPriv,
		Usage:           "sign",
		Status:          "active",
		ExpiresAt:       expiresAt,
		VerifyExpiresAt: expiresAt.Add(keyLifetime),
	})

	if err == nil {
		return nil
	}
	if fun.Is(err, fun.CodeConflict) {
		return nil
	}
	return err
}

func tryRotateProjectKeys(ctx context.Context, q *sqlc.Queries, encryptionKey []byte) error {
	dbe := database.NewErrorHandler("project key")
	projects, err := q.ListProjectsWithSigningKeys(ctx)
	err = dbe(err)
	if err != nil {
		return err
	}

	for _, projectID := range projects {
		var key sqlc.KeyPair
		key, err = q.GetActiveSigningKey(ctx, projectID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
				_ = createProjectKey(ctx, q, *projectID, encryptionKey)
				continue
			}
			return dbe(err)
		}

		if time.Until(key.ExpiresAt) > 24*time.Hour {
			continue
		}

		if err = q.RotateSigningKeys(ctx, projectID); err != nil {
			return dbe(err)
		}

		_ = createProjectKey(ctx, q, *projectID, encryptionKey)
	}

	return nil
}

func createProjectKey(ctx context.Context, q *sqlc.Queries, projectID uuid.UUID, encryptionKey []byte) error {
	dbe := database.NewErrorHandler("project key")
	pub, priv, err := crypto.GenerateEd25519()
	defer zero(priv)
	if err != nil {
		return err
	}

	encryptedPriv, err := crypto.Encrypt(priv, encryptionKey)
	if err != nil {
		return err
	}

	keyLifetime := errx.MustEnv("IDENTITY_X_KEY_LIFETIME", time.ParseDuration)
	kid := "project:" + projectID.String() + ":" + ulid.Make().String()
	expiresAt := time.Now().Add(keyLifetime)

	_, err = q.CreateKeyPair(ctx, sqlc.CreateKeyPairParams{
		Kid:             kid,
		ProjectID:       &projectID,
		KeyType:         "project",
		Algorithm:       "EdDSA",
		PublicKey:       pub,
		PrivateKey:      encryptedPriv,
		Usage:           "sign",
		Status:          "active",
		ExpiresAt:       expiresAt,
		VerifyExpiresAt: expiresAt.Add(keyLifetime),
	})

	// Rely on DB uniqueness for safety in concurrent rotations
	if err == nil {
		return nil
	}
	if fun.Is(err, fun.CodeConflict) {
		return nil
	}
	return dbe(err)
}

func queriesWithTx(ctx context.Context, q *sqlc.Queries) *sqlc.Queries {
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return q.WithTx(tx)
	}
	return q
}

func zero(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

func SetupCron(encryptionKey []byte, db *pgxpool.Pool, cfg Config) gocron.Scheduler {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Failed to create scheduler: %v", err)
	}

	txRunner := database.NewPGXTxRunner(db, telemetry.Log())
	rotateKeysJob(db, scheduler, txRunner, encryptionKey)
	sessionCleanupJob(db, scheduler, txRunner)
	tokenReuseCleanupJob(db, scheduler)

	go scheduler.Start()
	log.Println("Started the cron scheduler")
	return scheduler
}

func InitEncryption() []byte {
	encryptionKey, err := hex.DecodeString(os.Getenv("ENCRYPTION_KEY"))
	if err != nil {
		errx.Exit(err, "error decoding encryption key")
	}
	if len(encryptionKey) != 32 {
		errx.Exit(errors.New("encryption key size is not 32 bytes"), "Wrong key size")
	}
	return encryptionKey
}

func SetupAuthMiddlewares(
	sessions ports.SessionRepository,
	keys ports.KeysRepository,
	apiKeys ports.ApiKeyRepository,
	tracer trace.Tracer,
	issuer string,
) *middlewares.Middleware[*models.AccessClaims] {

	keyFunc := func(ctx context.Context, tokenStr string) (*models.AccessClaims, error) {
		ctx, span := tracer.Start(ctx, "Middleware.Auth.JWT")
		defer span.End()

		accessToken := &models.AccessClaims{}
		_, err := security.ParseJWTUnverified[*models.AccessClaims](tokenStr, accessToken)
		if err != nil {
			return nil, err
		}

		if accessToken.Sub.ProjectID != nil {
			span.SetAttributes(attribute.String("user.project_id", accessToken.Sub.ProjectID.String()))
		}

		keyPair, err := keys.GetActiveSigningKey(ctx, accessToken.Sub.ProjectID)
		if err != nil {
			return nil, err
		}

		accessToken, err = security.VerifyAccessToken(tokenStr, keyPair)
		if err != nil {
			return nil, err
		}

		if accessToken.Sub.ProjectID != nil && accessToken.Issuer != accessToken.Sub.ProjectID.String() {
			telemetry.DLog().Info("Project ID issuer branch", zap.String("issuer", accessToken.Issuer), zap.Any("project_id", accessToken.Sub.ProjectID))
			return nil, fun.ErrUnauthorized("access token has invalid issuer")
		} else if accessToken.Sub.ProjectID == nil && accessToken.Issuer != issuer {
			telemetry.DLog().Info("IDX native issuer branch", zap.String("passed issuer", issuer), zap.String("issuer", accessToken.Issuer), zap.Any("project_id", accessToken.Sub.ProjectID))
			return nil, fun.ErrUnauthorized("access token has invalid issuer")
		}

		sess, err := sessions.GetByFamilyID(ctx, accessToken.Sub.FamilyID)
		if err != nil {
			if fun.Is(err, fun.CodeNotFound) {
				return nil, fun.ErrUnauthorized("session not found or revoked")
			}
			return nil, err
		}

		if sess.SessionID != accessToken.Sub.SessionID {
			return nil, fun.ErrUnauthorized("token/session mismatch")
		}
		if sess.RevokedAt != nil {
			return nil, fun.ErrUnauthorized("session not found or revoked")
		}

		span.SetAttributes(
			attribute.String("user.type", accessToken.Sub.UserType),
			attribute.String("user.id", accessToken.Sub.ID.String()),
			attribute.String("user.session_id", accessToken.Sub.SessionID.String()),
		)

		return accessToken, nil
	}

	jwtHook := func(ctx context.Context, claims *models.AccessClaims) (context.Context, error) {
		principal, err := authz.NewPrincipal(claims)
		if err != nil {
			return ctx, err
		}
		return authz.WithPrincipal(ctx, principal), nil
	}

	apiKeyHook := func(ctx context.Context, rawKey string) (context.Context, error) {
		ctx, span := tracer.Start(ctx, "Middleware.Auth.APIKey")
		defer span.End()

		span.SetAttributes(attribute.String("auth.method", string(authz.AuthMethodApiKey)))

		if !strings.HasPrefix(rawKey, "gk_") {
			return ctx, fun.ErrUnauthorized("invalid api key shape")
		}

		parts := strings.SplitN(rawKey, "_", 3)
		if len(parts) != 3 {
			return ctx, fun.ErrUnauthorized("invalid api key shape")
		}

		projectID, err := uuid.Parse(parts[1])
		if err != nil {
			return ctx, fun.ErrUnauthorized("invalid api key shape")
		}

		keyData, err := apiKeys.GetByProjectID(ctx, projectID)
		if err != nil {
			if fun.Is(err, fun.CodeNotFound) {
				return ctx, fun.ErrUnauthorized("invalid api key")
			}
			return ctx, err
		}

		if err = crypto.VerifyBcryptSecret(keyData.KeyHash, parts[2]); err != nil {
			return ctx, fun.ErrUnauthorized("invalid api key")
		}

		return authz.WithPrincipal(ctx, &authz.Principal{
			UserID:    keyData.ClientID,
			ProjectID: &keyData.ProjectID,
			SessionID: nil,
			Method:    authz.AuthMethodApiKey,
		}), nil
	}

	return middlewares.New[*models.AccessClaims](keyFunc, jwtHook, apiKeyHook)
}
