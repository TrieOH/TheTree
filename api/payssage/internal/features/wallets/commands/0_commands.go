package commands

import (
	"payssage/ports"
)

type Commands struct {
	wallets ports.WalletRepo
	orgs    ports.OrganizationRepo
}

func NewCommands(
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
) *Commands {
	return &Commands{
		wallets: wallets,
		orgs:    orgs,
	}
}
