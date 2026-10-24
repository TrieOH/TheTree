package app

import (
	"context"
	"time"

	"lib/backend"
	"payssage/internal/authz"
	"payssage/internal/config"
	"payssage/internal/handlers"
	"payssage/internal/jobs"
	"payssage/internal/providers"
	"payssage/internal/providers/mercado_pago_provider"
	"payssage/internal/repos"
	"payssage/internal/services"
	"payssage/internal/sqlc"
	"payssage/ports"

	"github.com/go-chi/chi/v5"
	"resty.dev/v3"
)

// httpTimeout bounds Payssage's outgoing calls: provider APIs and tenant
// webhook deliveries.
const httpTimeout = 15 * time.Second

// Build wires Payssage on the kit (ADR-0014): the provider registry, the
// webhook jobs, the operations, and the strict server.
func Build(cfg config.Config) backend.BuildFunc {
	return func(_ context.Context, d *backend.Deps, r chi.Router) error {
		httpClient := resty.New().SetTimeout(httpTimeout)
		repos := repos.New(sqlc.New(d.DB))
		registry := paymentProviders(cfg, repos, httpClient)

		d.Jobs.Serve(jobs.Registry(jobs.Deps{
			Deliveries: repos.WebhookDeliveries,
			Events:     repos.WebhookEvents,
			Endpoints:  repos.WebhookEndpoints,
			HTTPClient: httpClient,
			Enqueuer:   d.Jobs,
		}), jobs.Periodic()...)

		ops := services.NewOperations(repos, authz.New(repos.Organizations, repos.Wallets), d.Jobs, d.IDX.Client, d.Tx, registry)

		chains, err := d.Chains(d.IDX.Auth)
		if err != nil {
			return err
		}
		mountStrict(r, handlers.NewServer(ops), chains)
		return nil
	}
}

// paymentProviders builds the provider registry. MercadoPago is the only
// provider and serves all three roles.
func paymentProviders(cfg config.Config, r *repos.Repos, httpClient *resty.Client) providers.PaymentProviders {
	mp := mercado_pago_provider.NewProvider(cfg.MercadoPagoConfig, r.Intents, r.Collectors, r.Sellers, r.Wallets, httpClient)
	return providers.PaymentProviders{
		OAuth:    map[providers.AvailableProviders]ports.OAuthProvider{providers.MercadoPagoProvider: mp},
		Payments: map[providers.AvailableProviders]ports.PaymentAbstractionLayer{providers.MercadoPagoProvider: mp},
		Webhooks: map[providers.AvailableProviders]ports.WebhookProvider{providers.MercadoPagoProvider: mp},
	}
}
