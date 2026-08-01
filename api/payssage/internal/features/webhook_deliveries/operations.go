package webhook_deliveries

import (
	"payssage/ports"
)

type Operations struct {
	deliveries ports.WebhookDeliveryRepo
	endpoints  ports.WebhookEndpointRepo
	wallets    ports.WalletRepo
	orgs       ports.OrganizationRepo
}

func NewOperations(
	deliveries ports.WebhookDeliveryRepo,
	endpoints ports.WebhookEndpointRepo,
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
) *Operations {
	return &Operations{
		deliveries: deliveries,
		endpoints:  endpoints,
		wallets:    wallets,
		orgs:       orgs,
	}
}
