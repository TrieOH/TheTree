package infrastructure

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/ports/inbounds"
	"database/sql"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Infra struct {
	DB      *sql.DB
	Queries *sqlc.Queries
	Tx      inbounds.TxRunner
	Logger  *zap.Logger
	Tracer  trace.Tracer
}

func NewInfra(db *sql.DB, queries *sqlc.Queries, tx inbounds.TxRunner, logger *zap.Logger, tracer trace.Tracer) Infra {
	return Infra{
		DB:      db,
		Queries: queries,
		Tx:      tx,
		Logger:  logger,
		Tracer:  tracer,
	}
}
