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

func NewQueries(deps ports.AuthnDeps) *Queries {
	return errx.MustProvide(&Queries{
		cryptoKeys: deps.CryptoKeys,
		logger:     deps.Logger,
		tracer:     deps.Tracer,
		tx:         deps.Tx,
	})
}
