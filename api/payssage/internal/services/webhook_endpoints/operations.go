package webhook_endpoints

import (
	"payssage/ports"
)

type Operations struct {
	endpoints ports.WebhookEndpointRepo
	wallets   ports.WalletRepo
	orgs      ports.OrganizationRepo
}

func NewOperations(
	endpoints ports.WebhookEndpointRepo,
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
) *Operations {
	return &Operations{
		endpoints: endpoints,
		wallets:   wallets,
		orgs:      orgs,
	}
}
