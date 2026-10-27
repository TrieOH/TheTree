package app

import (
	"context"

	"lib/backend"
	"lib/database"
	"univents/internal/authz"
	"univents/internal/config"
	"univents/internal/features/badges"
	"univents/internal/features/certifications"
	"univents/internal/features/checkouts"
	"univents/internal/features/editions"
	"univents/internal/features/events"
	"univents/internal/features/payments"
	"univents/internal/features/products"
	"univents/internal/features/programs"
	"univents/internal/features/purchases"
	"univents/internal/features/registrations"
	"univents/internal/features/signatures"
	"univents/internal/features/store"
	"univents/internal/features/ticket_types"
	"univents/internal/features/webhooks"
	"univents/internal/features/ws"
	"univents/internal/features/ws_tokens"
	"univents/internal/jobs"
	"univents/internal/sqlc"

	"github.com/go-chi/chi/v5"
)

// Build wires Univents on the kit (ADR-0014). Univents-only infrastructure
// — object storage, the Payssage client and its wallet check, the store
// notifier — is opened here rather than as modules: no other Backend uses
// it.
func Build(cfg config.Config) backend.BuildFunc {
	return func(ctx context.Context, d *backend.Deps, r chi.Router) error {
		storage, err := SetupObjectStorage(cfg)
		if err != nil {
			return err
		}
		pay := SetupPayssage(cfg)
		// Fail fast on the platform wallet: wrong/missing PAYSSAGE_WALLET_ID
		// or an unreachable Payssage must stop the boot, not surface at
		// checkout (D6).
		err = VerifyPayssageWallet(ctx, pay, cfg.PayssageWalletID)
		if err != nil {
			return err
		}
		// The store's LISTEN/NOTIFY bridge. Notify opens its own connection
		// per call, so nothing needs starting here.
		notifier := database.NewNotifier(cfg.ToDBConfig().DSN())

		q := sqlc.New(d.DB)
		eventsRepo := events.NewRepo(q)
		editionsRepo := editions.NewRepo(q)
		ticketTypesRepo := ticket_types.NewRepo(q)
		productsRepo := products.NewRepo(q)     // also the product purchases
		programsRepo := programs.NewRepo(q)     // also the occurrences and participations
		badgesRepo := badges.NewRepo(q)         // templates and emissions
		signaturesRepo := signatures.NewRepo(q) // also the signature requests
		certsRepo := certifications.NewRepo(q)
		registrationsRepo := registrations.NewRepo(q)
		purchasesRepo := purchases.NewRepo(q)
		wsTokensRepo := ws_tokens.NewRepo(q)
		az := authz.New(eventsRepo)
		actors := checkouts.NewSDKActorResolver(d.IDX.Client.Actors)

		d.Jobs.Serve(jobs.Registry(jobs.Deps{
			Certs:            certsRepo,
			Editions:         editionsRepo,
			Events:           eventsRepo,
			Purchases:        purchasesRepo,
			Registrations:    registrationsRepo,
			ProductPurchases: productsRepo,
			Participations:   programsRepo,
			TicketTypes:      ticketTypesRepo,
			Email:            d.Email,
			Notifier:         notifier,
			Tx:               d.Tx,
		}))

		// Badges and the WS hub are used by other features, so they are
		// built first.
		badgesOps := badges.New(badges.Deps{
			Repo: badgesRepo, Emissions: badgesRepo, Registrations: registrationsRepo, Editions: editionsRepo,
			Events: eventsRepo, Actors: actors, Email: d.Email, Authz: az,
		})
		wsOps := ws.New(ws.Deps{Tokens: wsTokensRepo, Purchases: purchasesRepo, Intents: pay, Notifier: notifier})

		srv := &server{
			eventHandlers: events.NewHandlers(events.New(events.Deps{
				Events: eventsRepo, Obj: storage, IDX: d.IDX.Client, Authz: az, Badges: badgesOps, Tx: d.Tx,
			})),
			editionHandlers: editions.NewHandlers(editions.New(editions.Deps{
				Events: eventsRepo, Editions: editionsRepo, Authz: az, Badges: badgesOps,
			})),
			ticketTypeHandlers: ticket_types.NewHandlers(ticket_types.New(ticket_types.Deps{
				Events: eventsRepo, Editions: editionsRepo, TicketTypes: ticketTypesRepo, Authz: az,
			})),
			productHandlers: products.NewHandlers(products.New(products.Deps{
				Events: eventsRepo, Editions: editionsRepo, Products: productsRepo, Authz: az, Tx: d.Tx,
			})),
			programHandlers: programs.NewHandlers(programs.New(programs.Deps{
				Events: eventsRepo, Editions: editionsRepo, Programs: programsRepo, Occurrences: programsRepo,
				Registrations: registrationsRepo, TicketTypes: ticketTypesRepo, Participations: programsRepo,
				Authz: az, Notifier: notifier, Tx: d.Tx,
			})),
			badgeHandlers: badges.NewHandlers(badgesOps),
			signatureHandlers: signatures.NewHandlers(signatures.New(signatures.Deps{
				Events: eventsRepo, Editions: editionsRepo, Signatures: signaturesRepo, Requests: signaturesRepo,
				Email: d.Email, HmacSecret: cfg.HmacSecret, Authz: az,
			})),
			certificationHandlers: certifications.NewHandlers(certifications.New(certifications.Deps{
				Events: eventsRepo, Editions: editionsRepo, Certs: certsRepo, Programs: programsRepo,
				Email: d.Email, Jobs: d.Jobs, Authz: az,
			})),
			paymentHandlers: payments.NewHandlers(payments.New(payments.Deps{
				Events: eventsRepo, Payssage: pay, Authz: az, WalletID: cfg.PayssageWalletID,
			})),
			webhookHandlers: webhooks.NewHandlers(webhooks.New(webhooks.Deps{
				Purchases: purchasesRepo, Registrations: registrationsRepo, ProductPurchases: productsRepo,
				Participations: programsRepo, Badges: badgesOps, Notifier: notifier, Jobs: d.Jobs, Tx: d.Tx,
				Secret: cfg.PayssageWebhookSecret,
			})),
			checkoutHandlers: checkouts.NewHandlers(checkouts.New(checkouts.Deps{
				Purchases: purchasesRepo, Editions: editionsRepo, Events: eventsRepo, TicketTypes: ticketTypesRepo,
				Products: productsRepo, Programs: programsRepo, Occurrences: programsRepo,
				Registrations: registrationsRepo, ProductPurchases: productsRepo, Participations: programsRepo,
				Badges: badgesOps, Notifier: notifier, Jobs: d.Jobs, Tx: d.Tx, Intents: pay, Payssage: pay,
				WalletID: cfg.PayssageWalletID, Tokens: wsOps, Actors: actors, Authz: az,
			})),
			purchaseHandlers: purchases.NewHandlers(purchases.New(purchases.Deps{Purchases: purchasesRepo})),
			wsHandlers:       ws.NewHandlers(wsOps),
			storeHandlers: store.NewHandlers(store.New(store.Deps{
				Avail: purchasesRepo, Editions: editionsRepo, Notifier: notifier,
			})),
		}

		chains, err := d.Chains(d.IDX.Auth)
		if err != nil {
			return err
		}
		mountStrict(r, srv, chains)
		mountRaw(r, srv)
		return nil
	}
}
