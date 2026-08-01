package sellers

import (
	"payssage/ports"
)

type Operations struct {
	sellers ports.SellerRepo
	wallets ports.WalletRepo
	orgs    ports.OrganizationRepo
}

func NewOperations(
	sellers ports.SellerRepo,
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
) *Operations {
	return &Operations{
		sellers: sellers,
		wallets: wallets,
		orgs:    orgs,
	}
}
