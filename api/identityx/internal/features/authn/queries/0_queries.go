package queries

import (
	"IdentityX/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Queries struct {
	projects   ports.ProjectRepo
	cryptoKeys ports.CryptoKeysRepo
	logger     *zap.Logger
	tracer     trace.Tracer
	tx         database.TxRunner
}

func NewQueries(
	projects ports.ProjectRepo,
	cryptoKeys ports.CryptoKeysRepo,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Queries {
	return &Queries{
		projects:   projects,
		cryptoKeys: cryptoKeys,
		logger:     logger,
		tracer:     tracer,
		tx:         tx,
	}
}
