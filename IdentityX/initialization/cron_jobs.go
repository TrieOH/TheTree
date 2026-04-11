package initialization

import (
	"GoAuth/internal/adapters/observability/logs"
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/ports/inbounds"
	"context"
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func rotateKeysJob(ctx context.Context, app *GoauthApp, txRunner inbounds.TxRunner) {

	scheduler := app.scheduler
	db := app.DB

	_, err := scheduler.NewJob(
		gocron.DurationJob(
			viper.GetDuration("ROTATE_KEYS_JOB_DURATION"),
		),
		gocron.NewTask(func(ctx context.Context, pool *pgxpool.Pool) {
			if err := txRunner.WithinTx(ctx, func(txCtx context.Context) error {
				q := sqlc.New(pool)
				q = queriesWithTx(txCtx, q)

				if err := tryRotateGoAuthKeys(txCtx, q); err != nil {
					return err
				}

				if err := tryRotateProjectKeys(txCtx, q); err != nil {
					return err
				}

				if err := q.RevokeExpiredRotatedKeys(txCtx); err != nil {
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
		}, db),
	)

	if err != nil {
		log.Fatalf("Couldn't create RotateKeys cron job: %v", err)
	}

	log.Println("Created RotateKeys cron job")
}

func sessionCleanupJob(ctx context.Context, app *GoauthApp, txRunner inbounds.TxRunner) {
	db := app.DB

	_, err := app.scheduler.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(0, 0, 0))),
		gocron.NewTask(func(ctx context.Context, pool *pgxpool.Pool) {
			ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()

			queries := sqlc.New(pool)

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
		}, db),
	)

	if err != nil {
		log.Fatalf("Couldn't create SessionCleanup cron job: %v", err)
	}

	log.Println("Created SessionCleanup cron job")
}

func tokenReuseCleanupJob(ctx context.Context, app *GoauthApp) {
	db := app.DB

	_, err := app.scheduler.NewJob(
		gocron.DurationJob(15*time.Minute),
		gocron.NewTask(func(ctx context.Context, pool *pgxpool.Pool) {
			ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()

			queries := sqlc.New(pool)

			if err := queries.DeleteExpiredTokenReuseListEntries(ctx); err != nil {
				logs.L().Error("Couldn't clear expired token reuse entries", zap.Error(err))
				return
			}
			logs.L().Info("Cleared expired token reuse entries")
		}, db),
	)

	if err != nil {
		log.Fatalf("Couldn't create TokenReuseCleanup cron job: %v", err)
	}

	log.Println("Created TokenReuseCleanup cron job")
}
