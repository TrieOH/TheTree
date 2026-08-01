package commands

import (
	"payssage/ports"
)

type Commands struct {
	endpoints ports.WebhookEndpointRepo
	wallets   ports.WalletRepo
	orgs      ports.OrganizationRepo
}

func NewCommands(
	endpoints ports.WebhookEndpointRepo,
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
) *Commands {
	return &Commands{
		endpoints: endpoints,
		wallets:   wallets,
		orgs:      orgs,
	}
}
