package intents

import (
	"payssage/ports"
)

type Operations struct {
	intents    ports.IntentRepo
	wallets    ports.WalletRepo
	orgs       ports.OrganizationRepo
	collectors ports.CollectorRepo
	sellers    ports.SellerRepo
}

func NewOperations(
	intents ports.IntentRepo,
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
	collectors ports.CollectorRepo,
	sellers ports.SellerRepo,
) *Operations {
	return &Operations{
		intents:    intents,
		wallets:    wallets,
		orgs:       orgs,
		collectors: collectors,
		sellers:    sellers,
	}
}
