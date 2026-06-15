package river

import (
	"context"
	"lib/errx"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivermigrate"
)

func Migrate(ctx context.Context, dbPool *pgxpool.Pool) {
	migrator, err := rivermigrate.New(riverpgxv5.New(dbPool), nil)
	if err != nil {
		errx.Exit(err, "failed to create migrator")
	}

	res, err := migrator.Migrate(ctx, rivermigrate.DirectionUp, nil)
	if err != nil {
		errx.Exit(err, "failed to migrate")
	}

	for _, v := range res.Versions {
		slog.Info("river migration applied", "version", v.Version)
	}
}
