package service

import (
	"payssage/internal/authz"

	"payssage/ports"
)

type Deps struct {
	Events  ports.WebhookEventRepo
	Wallets ports.WalletRepo
	Orgs    ports.OrganizationRepo
	Authz   *authz.Service
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }
