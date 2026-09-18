package app

import (
	"context"
	"fmt"
	libauthz "lib/authz"
	"lib/database"
	libriver "lib/river"
	"log/slog"
	"univents/internal/authz"
	"univents/internal/handlers"
	"univents/internal/repos"
	"univents/internal/services"
	certsJobs "univents/internal/services/certifications/jobs"
	checkoutsJobs "univents/internal/services/checkouts/jobs"
	"univents/internal/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"riverqueue.com/riverui"
)

// ── Init methods ──────────────────────────────────────────────────────────

func (app *Univents) initRepos() *repos.Repos {
	return repos.New(sqlc.New(app.db))
}

func (app *Univents) initOperations(r *repos.Repos, notifier *database.Notifier, riverClient *river.Client[pgx.Tx], tx database.TxRunner) *services.Operations {
	authzSvc := authz.New(r.Events)
	return services.NewOperations(r, authzSvc, app.objStorage, app.idxClient, app.emailClient, app.cfg.HmacSecret, app.payssage, app.cfg.PayssageWalletID, notifier, riverClient, tx, app.cfg.PayssageWebhookSecret)
}

func (app *Univents) initHandlers(ops *services.Operations) *handlers.Server {
	return handlers.NewServer(ops)
}

// initMiddlewares builds the auth primitives the spec-derived chains
// resolve against; construction stays per-backend, everything downstream
// (chain derivation, dispatch, fail-closed) is the Access-check and
// Harness modules' implementation.
func (app *Univents) initMiddlewares() libauthz.Primitives {
	authMW := app.setupAuthMiddlewares()

	return libauthz.Primitives{
		JWT:    authMW.JWT(),
		APIKey: authMW.APIKey(),
		Any:    authMW.AnyAuth(),
	}
}

func (app *Univents) initRiver(ctx context.Context, r *repos.Repos, notifier *database.Notifier, tx database.TxRunner) (*river.Client[pgx.Tx], *riverui.Handler, error) {
	libriver.Migrate(ctx, app.db)

	client := libriver.NewClient(app.db, libriver.NewWorkers(
		libriver.Register(certsJobs.NewGrantCertsWorker(r.Certs, r.Editions, r.Events, app.emailClient)),
		libriver.Register(certsJobs.NewGrantCertsForOccurrenceWorker(r.Certs, r.Editions, r.Events, app.emailClient)),
		libriver.Register(checkoutsJobs.NewExpirePurchaseWorker(r.Purchases, r.Registrations, r.Products, r.Programs, notifier, tx)),
		libriver.Register(checkoutsJobs.NewSendGiftEmailWorker(r.Registrations, r.Editions, r.Events, r.TicketTypes, app.emailClient)),
	), nil, nil)
	// TODO: schedule GrantCertsForEdition on edition end and GrantCertsForOccurrence on occurrence end

	err := client.Start(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("start river client: %w", err)
	}

	riverUIHandler, err := riverui.NewHandler(&riverui.HandlerOpts{
		DevMode:                  false,
		Endpoints:                riverui.NewEndpoints[pgx.Tx](client, nil),
		Logger:                   slog.Default(),
		Prefix:                   "/riverui",
		JobListHideArgsByDefault: true,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("create river ui handler: %w", err)
	}
	err = riverUIHandler.Start(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("start river ui handler: %w", err)
	}

	return client, riverUIHandler, nil
}
