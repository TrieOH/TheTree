package intents

import (
	"payssage/internal/authz"

	"payssage/ports"
)

type Operations struct {
	intents    ports.IntentRepo
	wallets    ports.WalletRepo
	orgs       ports.OrganizationRepo
	collectors ports.CollectorRepo
	sellers    ports.SellerRepo
	authz      *authz.Service
}

func NewOperations(
	intents ports.IntentRepo,
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
	collectors ports.CollectorRepo,
	sellers ports.SellerRepo,
	authz *authz.Service,
) *Operations {
	return &Operations{
		intents:    intents,
		wallets:    wallets,
		orgs:       orgs,
		collectors: collectors,
		sellers:    sellers,
		authz:      authz,
	}
}
