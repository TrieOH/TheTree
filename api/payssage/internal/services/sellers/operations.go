package sellers

import (
	"payssage/internal/authz"

	"payssage/ports"
)

type Operations struct {
	sellers ports.SellerRepo
	wallets ports.WalletRepo
	orgs    ports.OrganizationRepo
	authz   *authz.Service
}

func NewOperations(
	sellers ports.SellerRepo,
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
	authz *authz.Service,
) *Operations {
	return &Operations{
		sellers: sellers,
		wallets: wallets,
		orgs:    orgs,
		authz:   authz,
	}
}
