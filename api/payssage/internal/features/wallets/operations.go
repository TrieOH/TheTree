package wallets

import (
	"payssage/ports"
)

type Operations struct {
	wallets ports.WalletRepo
	orgs    ports.OrganizationRepo
}

func NewOperations(
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
) *Operations {
	return &Operations{
		wallets: wallets,
		orgs:    orgs,
	}
}
