package app

import (
	"IdentityX/internal/platform/database"
	"IdentityX/internal/platform/database/sqlc"
	"IdentityX/internal/shared/crypto"
	"IdentityX/internal/shared/errx"
	"context"
	"log"
	"time"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/fail/v3"
	"github.com/MintzyG/fail/v3/plugins/localization"
	"github.com/MintzyG/fail/v3/plugins/tracing/otel"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func SetupFail() {
	fail.AllowInternalLogs(true)
	fail.AllowStaticMutations(true, false)
	fail.AllowRuntimePanics(true)

	if err := fail.RegisterTranslator(&errx.HTTPTranslator{}); err != nil {
		log.Fatal(err)
	}

	fail.RegisterMapper(&errx.PGXMapper{})

	fail.SetLocalizer(localization.New())
	fail.SetDefaultLocale("en-US")

	tracerPlugin := otel.New(
		otel.WithTracerName("goauth-service"),
		otel.WithMode(otel.RecordSmart),
		otel.WithStackTrace(),
		otel.WithAttributePrefix("goauth.error"),
	)
	fail.SetTracer(tracerPlugin)
}

func SetupFUN() {
	module := viper.GetString("MODULE")
	if module == "" {
		module = "IdentityXAPI"
	}

	resp.SetConfig(resp.Config{
		MaxTraceSize:         50,
		ResponseSizeLimit:    10 * 1024 * 1024,
		MaxInterceptorAmount: 20,
		DefaultContentType:   "application/json",
		EnableSizeValidation: true,
		DefaultModule:        module,
		ErrorHandler:         errx.ErrToResp,
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
	var err error
	db, err := database.WaitForDB(30 * time.Second)
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}
	if err = database.RunMigrations(db, migrationPath); err != nil {
		log.Fatalf("Failed migrations: %v", err)
	}
	return db
}

func SetupRuntimeEnv(db *pgxpool.Pool) {
	queries := sqlc.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// First, rotate any expired keys to clear the way for new key creation
	if err := queries.RotateExpiredGoAuthKeys(ctx); err != nil {
		log.Printf("Warning: failed to rotate expired GoAuth keys: %v", err)
	}
	if err := queries.RotateExpiredProjectKeys(ctx); err != nil {
		log.Printf("Warning: failed to rotate expired project keys: %v", err)
	}

	// Also run the full key rotation logic to create new keys for projects without active keys
	if err := tryRotateGoAuthKeys(ctx, queries); err != nil {
		log.Printf("Warning: failed to rotate goauth keys: %v", err)
	}

	if err := tryRotateProjectKeys(ctx, queries); err != nil {
		log.Printf("Warning: failed to rotate project keys: %v", err)
	}

	_, err := queries.GetActiveSigningKeyForGoAuth(ctx)
	if err != nil {
		if fail.Is(fail.From(err), errx.SQLNotFound) {
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
			expiresAt := time.Now().Add(viper.GetDuration("GOAUTH_KEY_LIFETIME"))

			_, err = queries.CreateKeyPair(ctx, sqlc.CreateKeyPairParams{
				Kid:             kid,
				ProjectID:       nil,
				KeyType:         "goauth",
				Algorithm:       "EdDSA",
				PublicKey:       pub,
				PrivateKey:      encryptedPriv,
				Usage:           "sign",
				Status:          "active",
				ExpiresAt:       expiresAt,
				VerifyExpiresAt: expiresAt.Add(viper.GetDuration("GOAUTH_KEY_LIFETIME")),
			})

			fe := fail.From(err)

			if fe != nil {
				if errx.IsUniqueViolation(fe) {
					log.Println("GoAuth signing key already created by another instance")
				} else {
					log.Fatalf("failed to create GoAuth signing key: %v", fe.Error())
				}
			} else {
				log.Println("Created GoAuth signing key")
			}
		} else {
			log.Fatalf("failed checking GoAuth signing key: %v", err.Error())
		}
	}

	fe := fail.From(err)

	if fe != nil {
		if fail.Is(fe, errx.SCOPEOneGlobal) || errx.IsUniqueViolation(fe) {
			log.Println("GoAuth Global scope already created by another instance")
		} else {
			log.Fatalf("Failed to create GoAuth Global scope: %v", fe.Error())
		}
	} else {
		log.Println("Created GoAuth Global scope")
	}
}

func tryRotateGoAuthKeys(ctx context.Context, q *sqlc.Queries) error {
	key, err := q.GetActiveSigningKeyForGoAuth(ctx)
	if err != nil {
		if fail.Is(fail.From(err), errx.SQLNotFound) {
			// defensive: no signing key → create
			return createGoAuthKey(ctx, q)
		}
		return fail.From(err)
	}

	if time.Until(key.ExpiresAt) > 24*time.Hour {
		return nil
	}

	if err = q.RotateSigningKeysForGoAuth(ctx); err != nil {
		return fail.From(err)
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
	expiresAt := time.Now().Add(viper.GetDuration("GOAUTH_KEY_LIFETIME"))

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
		VerifyExpiresAt: expiresAt.Add(viper.GetDuration("GOAUTH_KEY_LIFETIME")),
	})

	if err == nil {
		return nil
	}
	fe := fail.From(err)
	if errx.IsUniqueViolation(fe) {
		return nil
	}
	return fe
}

func tryRotateProjectKeys(ctx context.Context, q *sqlc.Queries) error {
	projects, err := q.ListProjectsWithSigningKeys(ctx)
	fe := fail.From(err)
	if fe != nil {
		return fe
	}

	for _, projectID := range projects {
		var key sqlc.KeyPair
		key, err = q.GetActiveSigningKeyForProject(ctx, projectID)
		if err != nil {
			if fail.Is(fail.From(err), errx.SQLNotFound) {
				_ = createProjectKey(ctx, q, *projectID)
				continue
			}
			return fail.From(err)
		}

		if time.Until(key.ExpiresAt) > 24*time.Hour {
			continue
		}

		if err = q.RotateSigningKeysForProject(ctx, projectID); err != nil {
			return fail.From(err)
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
	expiresAt := time.Now().Add(viper.GetDuration("GOAUTH_KEY_LIFETIME"))

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
		VerifyExpiresAt: expiresAt.Add(viper.GetDuration("GOAUTH_KEY_LIFETIME")),
	})

	if err == nil {
		return nil
	}

	// Rely on DB uniqueness for safety in concurrent rotations
	if errx.IsUniqueViolation(fail.From(err)) {
		return nil
	}

	return fail.From(err)
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
