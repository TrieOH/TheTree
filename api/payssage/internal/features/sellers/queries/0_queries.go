package queries

import (
	"payssage/ports"
)

type Queries struct {
	sellers ports.SellerRepo
	wallets ports.WalletRepo
	orgs    ports.OrganizationRepo
}

func NewQueries(
	sellers ports.SellerRepo,
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
) *Queries {
	return &Queries{
		sellers: sellers,
		wallets: wallets,
		orgs:    orgs,
	}
}
