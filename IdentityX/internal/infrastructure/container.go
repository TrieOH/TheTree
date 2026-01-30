package infrastructure

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/ports/inbounds"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Infra struct {
	DB      *pgxpool.Pool
	Queries *sqlc.Queries
	Tx      inbounds.TxRunner
	Logger  *zap.Logger
	Tracer  trace.Tracer
}

func NewInfra(db *pgxpool.Pool, queries *sqlc.Queries, tx inbounds.TxRunner, logger *zap.Logger, tracer trace.Tracer) Infra {
	return Infra{
		DB:      db,
		Queries: queries,
		Tx:      tx,
		Logger:  logger,
		Tracer:  tracer,
	}
}
