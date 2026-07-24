package queries

import (
	"lib/database"
	"univents/ports"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Queries struct {
	signatures ports.SignatureRepo
	editions   ports.EditionRepo
	logger     *zap.Logger
	tracer     trace.Tracer
	tx         database.TxRunner
}

func NewQueries(
	signatures ports.SignatureRepo,
	editions ports.EditionRepo,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Queries {
	return &Queries{
		signatures: signatures,
		editions:   editions,
		logger:     logger,
		tracer:     tracer,
		tx:         tx,
	}
}
