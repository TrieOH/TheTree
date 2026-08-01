package queries

import (
	"payssage/ports"
)

type Queries struct {
	deliveries ports.WebhookDeliveryRepo
	endpoints  ports.WebhookEndpointRepo
	wallets    ports.WalletRepo
	orgs       ports.OrganizationRepo
}

func NewQueries(
	deliveries ports.WebhookDeliveryRepo,
	endpoints ports.WebhookEndpointRepo,
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
) *Queries {
	return &Queries{
		deliveries: deliveries,
		endpoints:  endpoints,
		wallets:    wallets,
		orgs:       orgs,
	}
}
