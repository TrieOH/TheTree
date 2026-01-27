package main

import (
	http2 "GoAuth/internal/adapters/http/handlers"
	"GoAuth/internal/adapters/observability/logs"
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/adapters/persistence/transactions"
	"GoAuth/internal/apierr"
	"GoAuth/internal/crypto"
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"GoAuth/internal/database"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/oklog/ulid/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var Port string
var DB *sql.DB
var scheduler gocron.Scheduler

// init initializes the application.
// It loads environment variables, sets up ED25519 keys, configures the response utility,
// connects to the database, runs migrations, sets the JWT master key, and schedules cron jobs.
func init() {
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
	if smtpUsername := viper.GetString("SMTP_USERNAME"); smtpUsername == "" {
		log.Fatalf("SMTP_USERNAME environment variable not set.")
	}
	if smtpPassword := viper.GetString("SMTP_PASSWORD"); smtpPassword == "" {
		log.Fatalf("SMTP_PASSWORD environment variable not set.")
	}
	if smtpFrom := viper.GetString("SMTP_FROM"); smtpFrom == "" {
		log.Fatalf("SMTP_FROM environment variable not set.")
	}

	env := viper.GetString("ENV")
	if env != "" && env != "production" {
		apierr.IncludeDebugCauses = true
	}

	Port = viper.GetString("PORT")
	if Port == "" {
		Port = "8080"
	}
	resp.SetConfig(resp.Config{
		MaxTraceSize:         50,
		ResponseSizeLimit:    10 * 1024 * 1024, // 10MB
		MaxInterceptorAmount: 20,
		DefaultContentType:   "application/json",
		EnableSizeValidation: true,
		DefaultModule:        "GoAuth-module",
		ErrorHandler:         apierr.ErrToResp,
	})

	var err error
	DB, err = database.WaitForDB(30 * time.Second)
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}

	if err := database.RunMigrations(DB, "./internal/database/migrations"); err != nil {
		log.Fatalf("Failed migrations: %v", err)
	}

	queries := sqlc.New(DB)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = queries.GetActiveSigningKeyForGoAuth(ctx)
	if err != nil {
		if apierr.IsNotFound(apierr.FromSQLC(err)) {
			// create new signing key
			pub, priv, err := crypto.GenerateEd25519()
			if err != nil {
				log.Fatalf("failed to generate GoAuth key: %v", err)
			}
			defer zero(priv)

			kid := "goauth:" + ulid.Make().String()
			expiresAt := time.Now().Add(7 * 24 * time.Hour)

			_, err = queries.CreateKeyPair(ctx, sqlc.CreateKeyPairParams{
				Kid:        kid,
				ProjectID:  nil,
				KeyType:    "goauth",
				Algorithm:  "EdDSA",
				PublicKey:  pub,
				PrivateKey: priv,
				Usage:      "sign",
				Status:     "active",
				ExpiresAt:  expiresAt,
			})

			if err != nil {
				// rely on DB uniqueness as safety net
				if apierr.IsUniqueViolation(err) {
					log.Println("GoAuth signing key already created by another instance")
				} else {
					log.Fatalf("failed to create GoAuth signing key: %v", err)
				}
			} else {
				log.Println("Created GoAuth signing key")
			}
		} else {
			log.Fatalf("failed checking GoAuth signing key: %v", err)
		}
	}

	// Create the scheduler
	scheduler, err = gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Failed to create scheduler: %v", err)
	}

	txRunner := transactions.NewTxRunner(DB)

	_, err = scheduler.NewJob(
		gocron.DurationJob(time.Hour),
		gocron.NewTask(func(ctx context.Context, db *sql.DB) {
			// Run all key rotations and cleanup inside a single transaction
			if err := txRunner.WithinTx(ctx, func(txCtx context.Context) error {
				q := sqlc.New(db)           // always create the base queries
				q = queriesWithTx(txCtx, q) // inject TX from context

				if err := tryRotateGoAuthKeys(txCtx, q); err != nil {
					return err // rollback automatically
				}

				if err := tryRotateProjectKeys(txCtx, q); err != nil {
					return err
				}

				if err := q.DeleteExpiredRevokedKeys(txCtx); err != nil {
					return err
				}

				logs.L().Info("Rotated GoAuth and projects keys")
				return nil
			}); err != nil {
				logs.L().Error("Scheduled key rotation failed, rolled back", zap.Error(err))
			}
		}, DB),
	)

	if err != nil {
		log.Fatalf("Couldn't create RotateKeys cron job: %v", err)
	} else {
		log.Println("Created RotateKeys cron job")
	}

	// SessionCleanup scheduled daily at 00:00
	_, err = scheduler.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(0, 0, 0))),
		gocron.NewTask(func(ctx context.Context, db *sql.DB) {
			ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()

			queries := sqlc.New(db)

			revoked, err := queries.RevokeExpiredSessions(ctx)
			if err != nil {
				logs.L().Error("Couldn't revoke expired sessions", zap.Error(err))
				return
			}
			logs.L().Info("Revoked expired sessions", zap.Int("count", len(revoked)))

			deleted, err := queries.DeleteRevokedSessions(ctx)
			if err != nil {
				logs.L().Error("Couldn't delete revoked sessions", zap.Error(err))
				return
			}
			logs.L().Info("Deleted revoked sessions", zap.Int("count", len(deleted)))
		}, DB),
	)

	if err != nil {
		log.Fatalf("Couldn't create SessionCleanup cron job: %v", err)
	} else {
		log.Println("Created SessionCleanup cron job")
	}

	// Start the scheduler in the background
	go scheduler.Start()
	log.Println("Started the cron scheduler")
}

func tryRotateGoAuthKeys(ctx context.Context, q *sqlc.Queries) error {
	key, err := q.GetActiveSigningKeyForGoAuth(ctx)
	if err != nil {
		if apierr.IsNotFound(err) {
			// defensive: no signing key → create
			return createGoAuthKey(ctx, q)
		}
		return err
	}

	if time.Until(key.ExpiresAt) > 24*time.Hour {
		return nil
	}

	if err := q.RotateSigningKeysForGoAuth(ctx); err != nil {
		return err
	}

	return createGoAuthKey(ctx, q)
}

func createGoAuthKey(ctx context.Context, q *sqlc.Queries) error {
	pub, priv, err := crypto.GenerateEd25519()
	defer zero(priv)
	if err != nil {
		return err
	}

	kid := "goauth:" + ulid.Make().String()
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	_, err = q.CreateKeyPair(ctx, sqlc.CreateKeyPairParams{
		Kid:        kid,
		ProjectID:  nil,
		KeyType:    "goauth",
		Algorithm:  "EdDSA",
		PublicKey:  pub,
		PrivateKey: priv,
		Usage:      "sign",
		Status:     "active",
		ExpiresAt:  expiresAt,
	})

	if apierr.IsUniqueViolation(err) {
		return nil
	}
	return err
}

func tryRotateProjectKeys(ctx context.Context, q *sqlc.Queries) error {
	projects, err := q.ListProjectsWithActiveSigningKeys(ctx)
	if err != nil {
		return err
	}

	for _, projectID := range projects {
		key, err := q.GetActiveSigningKeyForProject(ctx, projectID)
		if err != nil {
			if apierr.IsNotFound(err) {
				_ = createProjectKey(ctx, q, *projectID)
				continue
			}
			return err
		}

		if time.Until(key.ExpiresAt) > 24*time.Hour {
			continue
		}

		if err := q.RotateSigningKeysForProject(ctx, projectID); err != nil {
			return err
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

	kid := "project:" + projectID.String() + ":" + ulid.Make().String()
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	_, err = q.CreateKeyPair(ctx, sqlc.CreateKeyPairParams{
		Kid:        kid,
		ProjectID:  &projectID,
		KeyType:    "project",
		Algorithm:  "EdDSA",
		PublicKey:  pub,
		PrivateKey: priv,
		Usage:      "sign",
		Status:     "active",
		ExpiresAt:  expiresAt,
	})

	// Rely on DB uniqueness for safety in concurrent rotations
	if apierr.IsUniqueViolation(err) {
		return nil
	}

	return err
}

func queriesWithTx(ctx context.Context, q *sqlc.Queries) *sqlc.Queries {
	if tx, ok := ctx.Value(transactions.TxKeyValue).(*sql.Tx); ok && tx != nil {
		return q.WithTx(tx)
	}
	return q
}

func zero(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
