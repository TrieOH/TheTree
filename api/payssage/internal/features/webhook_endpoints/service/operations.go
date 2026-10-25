package service

import (
	"payssage/internal/authz"

	"payssage/ports"
)

type Deps struct {
	Endpoints ports.WebhookEndpointRepo
	Wallets   ports.WalletRepo
	Orgs      ports.OrganizationRepo
	Authz     *authz.Service
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }
