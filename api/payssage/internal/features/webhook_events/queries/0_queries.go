package queries

import (
	"payssage/ports"
)

type Queries struct {
	events  ports.WebhookEventRepo
	wallets ports.WalletRepo
	orgs    ports.OrganizationRepo
}

func NewQueries(
	events ports.WebhookEventRepo,
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
) *Queries {
	return &Queries{
		events:  events,
		wallets: wallets,
		orgs:    orgs,
	}
}
