package oauth

import (
	"payssage/ports"
)

type Operations struct {
	wallets    ports.WalletRepo
	orgs       ports.OrganizationRepo
	oauth      ports.OAuthStateRepo
	collectors ports.CollectorRepo
	sellers    ports.SellerRepo
}

func NewOperations(
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
	oauth ports.OAuthStateRepo,
	collectors ports.CollectorRepo,
	sellers ports.SellerRepo,
) *Operations {
	return &Operations{
		wallets:    wallets,
		orgs:       orgs,
		oauth:      oauth,
		collectors: collectors,
		sellers:    sellers,
	}
}
