package webhook_deliveries

import (
	"payssage/internal/authz"

	"payssage/ports"
)

type Operations struct {
	deliveries ports.WebhookDeliveryRepo
	endpoints  ports.WebhookEndpointRepo
	wallets    ports.WalletRepo
	orgs       ports.OrganizationRepo
	authz      *authz.Service
}

func NewOperations(
	deliveries ports.WebhookDeliveryRepo,
	endpoints ports.WebhookEndpointRepo,
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
	authz *authz.Service,
) *Operations {
	return &Operations{
		deliveries: deliveries,
		endpoints:  endpoints,
		wallets:    wallets,
		orgs:       orgs,
		authz:      authz,
	}
}
