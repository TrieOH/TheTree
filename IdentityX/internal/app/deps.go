package app

import (
	"IdentityX/internal/platform/database"
	"IdentityX/internal/platform/database/sqlc"
	"IdentityX/internal/shared/crypto"
	"IdentityX/internal/shared/errx"
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"
	"unicode"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
	"github.com/go-chi/chi/v5"
	"github.com/go-co-op/gocron/v2"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
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
		errx.Must(err, "failed to register uuid7 validator")
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
		errx.Must(err, "failed to register passwd validator")
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
	module := viper.GetString("MODULE")
	if module == "" {
		module = "IdentityXAPI"
	}

	fun.SetConfig(fun.Config{
		MaxTraceSize:         50,
		ResponseSizeLimit:    10 * 1024 * 1024,
		MaxInterceptorAmount: 20,
		DefaultContentType:   "application/json",
		EnableSizeValidation: true,
		DefaultModule:        module,
	})

	v := SetupValidator()
	bind.SetValidator(v)
	fun.SetPathParamFunc(func(r *http.Request, key string) string {
		return chi.URLParam(r, key)
	})
}

func SetupRedis(timeout time.Duration) *redis.Client {
	rdb, err := database.WaitForRedis(timeout)
	if err != nil {
		log.Fatalf("Failed to connect Redis: %v", err)
	}
	return rdb
}

func SetupDB(migrationPath string) *pgxpool.Pool {
	db, err := database.WaitForDB(30 * time.Second)
	if err != nil {
		errx.Must(err, "Failed to connect DB")
	}
	if err = database.RunMigrations(db, migrationPath); err != nil {
		log.Fatalf("Failed migrations: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	errx.Must(errx.ValidateConstraintRegistry(ctx, db), "unregistered constraints found")
	return db
}

func SetupRuntimeEnv(db *pgxpool.Pool) {
	q := sqlc.New(db)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// First, rotate any expired keys to clear the way for new key creation
	if err := q.RotateExpiredKeys(ctx); err != nil {
		log.Printf("Warning: failed to rotate expired signing keys: %v", err)
	}

	// Also run the full key rotation logic to create new keys for projects without active keys
	if err := tryRotateGoAuthKeys(ctx, q); err != nil {
		log.Printf("Warning: failed to rotate goauth keys: %v", err)
	}

	if err := tryRotateProjectKeys(ctx, q); err != nil {
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
			encryptedPriv, err = crypto.Encrypt(priv)
			if err != nil {
				log.Fatalf("failed to encrypt GoAuth key: %v", err)
			}

			kid := "goauth:" + ulid.Make().String()
			expiresAt := time.Now().Add(viper.GetDuration("IDENTITY_X_KEY_LIFETIME"))

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
				VerifyExpiresAt: expiresAt.Add(viper.GetDuration("IDENTITY_X_KEY_LIFETIME")),
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

func tryRotateGoAuthKeys(ctx context.Context, q *sqlc.Queries) error {
	key, err := q.GetActiveSigningKey(ctx, nil)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			// defensive: no signing key → create
			return createGoAuthKey(ctx, q)
		}
		return errx.DB(err, "Signing Key")
	}

	if time.Until(key.ExpiresAt) > 24*time.Hour {
		return nil
	}

	if err = q.RotateSigningKeys(ctx, nil); err != nil {
		return errx.DB(err, "Signing Key")
	}

	return createGoAuthKey(ctx, q)
}

func createGoAuthKey(ctx context.Context, q *sqlc.Queries) error {
	pub, priv, err := crypto.GenerateEd25519()
	defer zero(priv)
	if err != nil {
		return err
	}

	encryptedPriv, err := crypto.Encrypt(priv)
	if err != nil {
		return err
	}

	kid := "goauth:" + ulid.Make().String()
	expiresAt := time.Now().Add(viper.GetDuration("IDENTITY_X_KEY_LIFETIME"))

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
		VerifyExpiresAt: expiresAt.Add(viper.GetDuration("IDENTITY_X_KEY_LIFETIME")),
	})

	if err == nil {
		return nil
	}
	if fun.Is(err, fun.CodeConflict) {
		return nil
	}
	return err
}

func tryRotateProjectKeys(ctx context.Context, q *sqlc.Queries) error {
	projects, err := q.ListProjectsWithSigningKeys(ctx)
	err = errx.DB(err, "Signing Key")
	if err != nil {
		return err
	}

	for _, projectID := range projects {
		var key sqlc.KeyPair
		key, err = q.GetActiveSigningKey(ctx, projectID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
				_ = createProjectKey(ctx, q, *projectID)
				continue
			}
			return errx.DB(err, "project Keys")
		}

		if time.Until(key.ExpiresAt) > 24*time.Hour {
			continue
		}

		if err = q.RotateSigningKeys(ctx, projectID); err != nil {
			return errx.DB(err, "project Keys")
		}

		_ = createProjectKey(ctx, q, *projectID)
	}

	return nil
}

func createProjectKey(ctx context.Context, q *sqlc.Queries, projectID uuid.UUID) error {
	pub, priv, err := crypto.GenerateEd25519()
	defer zero(priv)
	if err != nil {
		return err
	}

	encryptedPriv, err := crypto.Encrypt(priv)
	if err != nil {
		return err
	}

	kid := "project:" + projectID.String() + ":" + ulid.Make().String()
	expiresAt := time.Now().Add(viper.GetDuration("IDENTITY_X_KEY_LIFETIME"))

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
		VerifyExpiresAt: expiresAt.Add(viper.GetDuration("IDENTITY_X_KEY_LIFETIME")),
	})

	// Rely on DB uniqueness for safety in concurrent rotations
	if err == nil {
		return nil
	}
	if fun.Is(err, fun.CodeConflict) {
		return nil
	}
	return errx.DB(err, "project Keys")
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

func SetupCron(db *pgxpool.Pool) gocron.Scheduler {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Failed to create scheduler: %v", err)
	}

	txRunner := database.NewPGTxRunner(db)
	rotateKeysJob(db, scheduler, txRunner)
	sessionCleanupJob(db, scheduler, txRunner)
	tokenReuseCleanupJob(db, scheduler)

	go scheduler.Start()
	log.Println("Started the cron scheduler")
	return scheduler
}
