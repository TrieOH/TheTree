package queries

import (
	"IdentityX/ports"
	"lib/database"
	"lib/errx"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Queries struct {
	cryptoKeys ports.CryptoKeysRepo
	logger     *zap.Logger
	tracer     trace.Tracer
	tx         database.TxRunner
}

func NewQueries(
	cryptoKeys ports.CryptoKeysRepo,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Queries {
	return errx.MustProvide(&Queries{
		cryptoKeys: cryptoKeys,
		logger:     logger,
		tracer:     tracer,
		tx:         tx,
	})
}
