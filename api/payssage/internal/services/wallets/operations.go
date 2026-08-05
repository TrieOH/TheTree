package wallets

import (
	"payssage/internal/authz"

	"payssage/ports"
)

type Operations struct {
	wallets ports.WalletRepo
	orgs    ports.OrganizationRepo
	authz   *authz.Service
}

func NewOperations(
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
	authz *authz.Service,
) *Operations {
	return &Operations{
		wallets: wallets,
		orgs:    orgs,
		authz:   authz,
	}
}
