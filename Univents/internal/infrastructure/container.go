package infrastructure

import (
	"univents/internal/adapters/persistence/sqlc"
	"univents/internal/ports/inbounds"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Infra struct {
	DB      *pgxpool.Pool
	Queries *sqlc.Queries
	Tx      inbounds.TxRunner
	Logger  *zap.Logger
	Tracer  trace.Tracer
	Redis   *redis.Client
}

func NewInfra(db *pgxpool.Pool, queries *sqlc.Queries, tx inbounds.TxRunner, logger *zap.Logger, tracer trace.Tracer, rdb *redis.Client) Infra {
	return Infra{
		DB:      db,
		Queries: queries,
		Tx:      tx,
		Logger:  logger,
		Tracer:  tracer,
		Redis:   rdb,
	}
}
