package queries

import (
	"payssage/ports"
)

type Queries struct {
	wallets ports.WalletRepo
	orgs    ports.OrganizationRepo
}

func NewQueries(
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
) *Queries {
	return &Queries{
		wallets: wallets,
		orgs:    orgs,
	}
}
