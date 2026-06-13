package app

import (
	"context"
	"log"
	"time"

	"IdentityX/internal/database/sqlc"
	"lib/database"
	"lib/errx"
	"lib/telemetry"

	"github.com/go-co-op/gocron/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func rotateKeysJob(db *pgxpool.Pool, scheduler gocron.Scheduler, txRunner database.TxRunner, encryptionKey []byte) {
	rotateJobDuration := errx.MustEnv("ROTATE_KEYS_JOB_DURATION", time.ParseDuration)
	_, err := scheduler.NewJob(
		gocron.DurationJob(rotateJobDuration),
		gocron.NewTask(func(ctx context.Context, pool *pgxpool.Pool) {
			if err := txRunner.WithinTx(ctx, func(txCtx context.Context) error {
				q := sqlc.New(pool)
				q = queriesWithTx(txCtx, q)
				if err := tryRotateIDXKeys(txCtx, q, encryptionKey); err != nil {
					return err
				}
				if err := tryRotateProjectKeys(txCtx, q, encryptionKey); err != nil {
					return err
				}
				if err := q.RevokeExpiredRotatedKeys(txCtx); err != nil {
					return err
				}
				if err := q.DeleteExpiredRevokedKeys(txCtx); err != nil {
					return err
				}
				telemetry.Log().Info("Rotated GoAuth and projects keys")
				return nil
			}); err != nil {
				telemetry.Log().Error("Scheduled key rotation failed, rolled back", zap.Error(err))
			}
		}, db),
	)
	if err != nil {
		log.Fatalf("Couldn't create RotateKeys cron job: %v", err)
	}
	log.Println("Created RotateKeys cron job")
}

func sessionCleanupJob(db *pgxpool.Pool, scheduler gocron.Scheduler, txRunner database.TxRunner) {
	_, err := scheduler.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(0, 0, 0))),
		gocron.NewTask(func(ctx context.Context, pool *pgxpool.Pool) {
			ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()
			queries := sqlc.New(pool)
			revoked, err := queries.RevokeExpiredSessions(ctx)
			if err != nil {
				telemetry.Log().Error("Couldn't revoke expired sessions", zap.Error(err))
				return
			}
			telemetry.Log().Info("Revoked expired sessions", zap.Int("count", len(revoked)))
			deleted, err := queries.DeleteRevokedSessions(ctx)
			if err != nil {
				telemetry.Log().Error("Couldn't delete revoked sessions", zap.Error(err))
				return
			}
			telemetry.Log().Info("Deleted revoked sessions", zap.Int("count", len(deleted)))
		}, db),
	)
	if err != nil {
		log.Fatalf("Couldn't create SessionCleanup cron job: %v", err)
	}
	log.Println("Created SessionCleanup cron job")
}

func tokenReuseCleanupJob(db *pgxpool.Pool, scheduler gocron.Scheduler) {
	_, err := scheduler.NewJob(
		gocron.DurationJob(15*time.Minute),
		gocron.NewTask(func(ctx context.Context, pool *pgxpool.Pool) {
			ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()
			queries := sqlc.New(pool)
			if err := queries.DeleteExpiredTokenReuseListEntries(ctx); err != nil {
				telemetry.Log().Error("Couldn't clear expired token reuse entries", zap.Error(err))
				return
			}
			telemetry.Log().Info("Cleared expired token reuse entries")
		}, db),
	)
	if err != nil {
		log.Fatalf("Couldn't create TokenReuseCleanup cron job: %v", err)
	}
	log.Println("Created TokenReuseCleanup cron job")
}
