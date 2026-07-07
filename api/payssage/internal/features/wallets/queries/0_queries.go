package queries

import (
	"lib/database"
	"payssage/ports"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Queries struct {
	wallets ports.WalletRepo
	orgs    ports.OrganizationRepo
	logger  *zap.Logger
	tracer  trace.Tracer
	tx      database.TxRunner
}

func NewQueries(
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Queries {
	return &Queries{
		wallets: wallets,
		orgs:    orgs,
		logger:  logger,
		tracer:  tracer,
		tx:      tx,
	}
}
