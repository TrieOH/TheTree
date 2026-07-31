package queries

import (
	"payssage/ports"
)

type Queries struct {
	endpoints ports.WebhookEndpointRepo
	wallets   ports.WalletRepo
	orgs      ports.OrganizationRepo
}

func NewQueries(
	endpoints ports.WebhookEndpointRepo,
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
) *Queries {
	return &Queries{
		endpoints: endpoints,
		wallets:   wallets,
		orgs:      orgs,
	}
}
