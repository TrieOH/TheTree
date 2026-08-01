package webhook_events

import (
	"payssage/ports"
)

type Operations struct {
	events  ports.WebhookEventRepo
	wallets ports.WalletRepo
	orgs    ports.OrganizationRepo
}

func NewOperations(
	events ports.WebhookEventRepo,
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
) *Operations {
	return &Operations{
		events:  events,
		wallets: wallets,
		orgs:    orgs,
	}
}
