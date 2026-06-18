package queries

import (
	"IdentityX/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Queries struct {
	projects ports.ProjectRepo
	actors   ports.ActorRepo
	logger   *zap.Logger
	tracer   trace.Tracer
	tx       database.TxRunner
}

func NewQueries(
	actors ports.ActorRepo,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Queries {
	return &Queries{
		actors: actors,
		logger: logger,
		tracer: tracer,
		tx:     tx,
	}
}
