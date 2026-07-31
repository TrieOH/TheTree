package commands

import (
	"payssage/ports"
)

type Commands struct {
	intents    ports.IntentRepo
	wallets    ports.WalletRepo
	orgs       ports.OrganizationRepo
	collectors ports.CollectorRepo
	sellers    ports.SellerRepo
}

func NewCommands(
	intents ports.IntentRepo,
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
	collectors ports.CollectorRepo,
	sellers ports.SellerRepo,
) *Commands {
	return &Commands{
		intents:    intents,
		wallets:    wallets,
		orgs:       orgs,
		collectors: collectors,
		sellers:    sellers,
	}
}
