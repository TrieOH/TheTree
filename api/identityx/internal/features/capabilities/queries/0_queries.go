package queries

import (
	"IdentityX/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Queries struct {
	capabilities ports.CapabilityRepo
	projects     ports.ProjectRepo
	logger       *zap.Logger
	tracer       trace.Tracer
	tx           database.TxRunner
}

func NewQueries(
	capabilities ports.CapabilityRepo,
	projects ports.ProjectRepo,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Queries {
	return &Queries{
		capabilities: capabilities,
		projects:     projects,
		logger:       logger,
		tracer:       tracer,
		tx:           tx,
	}
}
