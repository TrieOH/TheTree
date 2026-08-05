package oauth

import (
	"payssage/internal/authz"

	"payssage/ports"
)

type Operations struct {
	wallets    ports.WalletRepo
	orgs       ports.OrganizationRepo
	oauth      ports.OAuthStateRepo
	collectors ports.CollectorRepo
	sellers    ports.SellerRepo
	authz      *authz.Service
}

func NewOperations(
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
	oauth ports.OAuthStateRepo,
	collectors ports.CollectorRepo,
	sellers ports.SellerRepo,
	authz *authz.Service,
) *Operations {
	return &Operations{
		wallets:    wallets,
		orgs:       orgs,
		oauth:      oauth,
		collectors: collectors,
		sellers:    sellers,
		authz:      authz,
	}
}
