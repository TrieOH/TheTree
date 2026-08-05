package webhook_events

import (
	"payssage/internal/authz"

	"payssage/ports"
)

type Operations struct {
	events  ports.WebhookEventRepo
	wallets ports.WalletRepo
	orgs    ports.OrganizationRepo
	authz   *authz.Service
}

func NewOperations(
	events ports.WebhookEventRepo,
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
	authz *authz.Service,
) *Operations {
	return &Operations{
		events:  events,
		wallets: wallets,
		orgs:    orgs,
		authz:   authz,
	}
}
