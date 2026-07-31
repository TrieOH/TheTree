package queries

import (
	"payssage/ports"
)

type Queries struct {
	intents ports.IntentRepo
	wallets ports.WalletRepo
	orgs    ports.OrganizationRepo
}

func NewQueries(
	intents ports.IntentRepo,
	wallets ports.WalletRepo,
	orgs ports.OrganizationRepo,
) *Queries {
	return &Queries{
		intents: intents,
		wallets: wallets,
		orgs:    orgs,
	}
}
