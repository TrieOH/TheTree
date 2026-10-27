package service

import (
	"payssage/internal/authz"

	"payssage/internal/providers"
	"payssage/ports"
)

type Deps struct {
	Intents    ports.IntentRepo
	Wallets    ports.WalletRepo
	Orgs       ports.OrganizationRepo
	Collectors ports.CollectorRepo
	Sellers    ports.SellerRepo
	Authz      *authz.Service
	Registry   providers.PaymentProviders
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }
