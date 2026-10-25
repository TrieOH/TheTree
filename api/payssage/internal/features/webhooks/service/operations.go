package service

import (
	"payssage/internal/authz"

	"payssage/internal/providers"
	"payssage/ports"

	"lib/database"
	libjobs "lib/jobs"
)

type Deps struct {
	Jobs       libjobs.Enqueuer
	Events     ports.WebhookEventRepo
	Endpoints  ports.WebhookEndpointRepo
	Deliveries ports.WebhookDeliveryRepo
	Authz      *authz.Service
	Registry   providers.PaymentProviders
	// Tx opens the transactions multi-step writes run in; threaded
	// from boot, no package-level runner.
	Tx database.TxRunner
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }
