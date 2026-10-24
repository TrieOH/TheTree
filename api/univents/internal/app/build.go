package app

import (
	"context"

	"lib/backend"
	"lib/database"
	"univents/internal/authz"
	"univents/internal/config"
	"univents/internal/handlers"
	"univents/internal/jobs"
	"univents/internal/repos"
	"univents/internal/services"
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

		repos := repos.New(sqlc.New(d.DB))

		d.Jobs.Serve(jobs.Registry(jobs.Deps{
			Certs:            repos.Certs,
			Editions:         repos.Editions,
			Events:           repos.Events,
			Purchases:        repos.Purchases,
			Registrations:    repos.Registrations,
			ProductPurchases: repos.Products,
			Participations:   repos.Programs,
			TicketTypes:      repos.TicketTypes,
			Email:            d.Email,
			Notifier:         notifier,
			Tx:               d.Tx,
		}))

		ops := services.NewOperations(repos, authz.New(repos.Events), storage, d.IDX.Client, d.Email,
			cfg.HmacSecret, pay, cfg.PayssageWalletID, notifier, d.Jobs, d.Tx, cfg.PayssageWebhookSecret)
		h := handlers.NewServer(ops)

		chains, err := d.Chains(d.IDX.Auth)
		if err != nil {
			return err
		}
		mountStrict(r, h, chains)
		mountRaw(r, h)
		return nil
	}
}
