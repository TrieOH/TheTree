package app

import (
	"context"
	"time"

	"lib/backend"
	"payssage/internal/authz"
	"payssage/internal/config"
	"payssage/internal/features/collectors"
	"payssage/internal/features/intents"
	"payssage/internal/features/oauth"
	"payssage/internal/features/orgs"
	"payssage/internal/features/sellers"
	"payssage/internal/features/wallets"
	"payssage/internal/features/webhook_deliveries"
	"payssage/internal/features/webhook_endpoints"
	"payssage/internal/features/webhook_events"
	"payssage/internal/features/webhooks"
	"payssage/internal/jobs"
	"payssage/internal/openapi"
	"payssage/internal/providers"
	"payssage/internal/providers/mercado_pago_provider"
	"payssage/internal/sqlc"
	"payssage/ports"

	"github.com/go-chi/chi/v5"
	"resty.dev/v3"
)

// httpTimeout bounds Payssage's outgoing calls: provider APIs and tenant
// webhook deliveries.
const httpTimeout = 15 * time.Second

// Build wires Payssage on the kit (ADR-0014): the repos, the provider
// registry, the webhook jobs, each feature's operations (ADR-0015), and the
// strict server.
func Build(cfg config.Config) backend.BuildFunc {
	return func(_ context.Context, d *backend.Deps, r chi.Router) error {
		httpClient := resty.New().SetTimeout(httpTimeout)

		q := sqlc.New(d.DB)
		orgsRepo := orgs.NewRepo(q)
		walletsRepo := wallets.NewRepo(q)
		oauthRepo := oauth.NewRepo(q)
		collectorsRepo := collectors.NewRepo(q)
		sellersRepo := sellers.NewRepo(q)
		intentsRepo := intents.NewRepo(q)
		endpointsRepo := webhook_endpoints.NewRepo(q)
		deliveriesRepo := webhook_deliveries.NewRepo(q)
		eventsRepo := webhook_events.NewRepo(q)
		az := authz.New(orgsRepo, walletsRepo)

		mp := mercado_pago_provider.NewProvider(cfg.MercadoPagoConfig, intentsRepo, collectorsRepo, sellersRepo, walletsRepo, httpClient)
		// MercadoPago is the only provider and serves all three roles.
		registry := providers.PaymentProviders{
			OAuth:    map[providers.AvailableProviders]ports.OAuthProvider{providers.MercadoPagoProvider: mp},
			Payments: map[providers.AvailableProviders]ports.PaymentAbstractionLayer{providers.MercadoPagoProvider: mp},
			Webhooks: map[providers.AvailableProviders]ports.WebhookProvider{providers.MercadoPagoProvider: mp},
		}

		d.Jobs.Serve(jobs.Registry(jobs.Deps{
			Deliveries: deliveriesRepo,
			Events:     eventsRepo,
			Endpoints:  endpointsRepo,
			HTTPClient: httpClient,
			Enqueuer:   d.Jobs,
		}), jobs.Periodic()...)

		srv := &server{
			orgHandlers: orgs.NewHandlers(orgs.New(orgs.Deps{
				Orgs: orgsRepo, IDX: d.IDX.Client, Authz: az, Tx: d.Tx,
			})),
			walletHandlers: wallets.NewHandlers(wallets.New(wallets.Deps{
				Wallets: walletsRepo, Orgs: orgsRepo, Authz: az,
			})),
			oauthHandlers: oauth.NewHandlers(oauth.New(oauth.Deps{
				Wallets: walletsRepo, Orgs: orgsRepo, OAuth: oauthRepo, Collectors: collectorsRepo,
				Sellers: sellersRepo, Authz: az, Registry: registry,
			})),
			collectorHandlers: collectors.NewHandlers(collectors.New(collectors.Deps{
				Collectors: collectorsRepo, Orgs: orgsRepo, Authz: az,
			})),
			sellerHandlers: sellers.NewHandlers(sellers.New(sellers.Deps{
				Sellers: sellersRepo, Wallets: walletsRepo, Orgs: orgsRepo, Authz: az,
			})),
			intentHandlers: intents.NewHandlers(intents.New(intents.Deps{
				Intents: intentsRepo, Wallets: walletsRepo, Orgs: orgsRepo, Collectors: collectorsRepo,
				Sellers: sellersRepo, Authz: az, Registry: registry,
			})),
			webhookHandlers: webhooks.NewHandlers(webhooks.New(webhooks.Deps{
				Jobs: d.Jobs, Events: eventsRepo, Endpoints: endpointsRepo, Deliveries: deliveriesRepo,
				Authz: az, Registry: registry, Tx: d.Tx,
			})),
			webhookEndpointHandlers: webhook_endpoints.NewHandlers(webhook_endpoints.New(webhook_endpoints.Deps{
				Endpoints: endpointsRepo, Wallets: walletsRepo, Orgs: orgsRepo, Authz: az,
			})),
			webhookEventHandlers: webhook_events.NewHandlers(webhook_events.New(webhook_events.Deps{
				Events: eventsRepo, Wallets: walletsRepo, Orgs: orgsRepo, Authz: az,
			})),
			webhookDeliveryHandlers: webhook_deliveries.NewHandlers(webhook_deliveries.New(webhook_deliveries.Deps{
				Deliveries: deliveriesRepo, Endpoints: endpointsRepo, Wallets: walletsRepo, Orgs: orgsRepo, Authz: az,
			})),
		}

		chains, err := d.Chains(d.IDX.Auth)
		if err != nil {
			return err
		}
		// The raw-request capture runs before the strict server decodes the
		// body, so the provider webhook receive verifies signatures against
		// the exact bytes. Path-scoped to /webhooks/ inside the middleware.
		openapi.Mount(r, srv, chains, webhooks.RawRequestMiddleware)
		return nil
	}
}
