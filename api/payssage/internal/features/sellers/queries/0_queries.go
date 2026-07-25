package queries

import (
	"lib/database"
	"payssage/ports"

	"go.opentelemetry.io/otel/trace"
)

type Queries struct {
	sellers ports.SellerRepo
	wallets ports.WalletRepo
	orgs    ports.OrganizationRepo
	tracer  trace.Tracer
	tx      database.TxRunner
}

func NewQueries(
	sellers ports.SellerRepo,
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
	tracer trace.Tracer,
	tx database.TxRunner,
) *Queries {
	return &Queries{
		sellers: sellers,
		wallets: wallets,
		orgs:    orgs,
		tracer:  tracer,
		tx:      tx,
	}
}
