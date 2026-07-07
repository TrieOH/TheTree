package commands

import (
	"lib/database"
	"payssage/ports"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Commands struct {
	wallets ports.WalletRepo
	orgs    ports.OrganizationRepo
	logger  *zap.Logger
	tracer  trace.Tracer
	tx      database.TxRunner
}

func NewCommands(
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Commands {
	return &Commands{
		wallets: wallets,
		orgs:    orgs,
		logger:  logger,
		tracer:  tracer,
		tx:      tx,
	}
}
