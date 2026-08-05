package webhook_endpoints

import (
	"payssage/internal/authz"

	"payssage/ports"
)

type Operations struct {
	endpoints ports.WebhookEndpointRepo
	wallets   ports.WalletRepo
	orgs      ports.OrganizationRepo
	authz     *authz.Service
}

func NewOperations(
	endpoints ports.WebhookEndpointRepo,
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
	authz *authz.Service,
) *Operations {
	return &Operations{
		endpoints: endpoints,
		wallets:   wallets,
		orgs:      orgs,
		authz:     authz,
	}
}
