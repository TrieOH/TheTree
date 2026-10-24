package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	libauthz "lib/authz"
	"lib/database"
	libjobs "lib/jobs"
	"lib/jobs/httpjobs"
	"lib/jobs/inprocess"
	"univents/internal/authz"
	"univents/internal/handlers"
	"univents/internal/jobs"
	"univents/internal/repos"
	"univents/internal/services"
	"univents/internal/sqlc"
)

// ── Init methods ──────────────────────────────────────────────────────────

func (app *Univents) initRepos() *repos.Repos {
	return repos.New(sqlc.New(app.db))
}

func (app *Univents) initOperations(r *repos.Repos, notifier *database.Notifier, enqueuer libjobs.DelayedEnqueuer, tx database.TxRunner) *services.Operations {
	authzSvc := authz.New(r.Events)
	return services.NewOperations(r, authzSvc, app.objStorage, app.idxClient, app.emailClient, app.cfg.HmacSecret, app.payssage, app.cfg.PayssageWalletID, notifier, enqueuer, tx, app.cfg.PayssageWebhookSecret)
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

// jobsDrainTimeout bounds how long shutdown waits for queued jobs.
const jobsDrainTimeout = 15 * time.Second

// initJobs builds the job registry, picks the carrier, and returns the
// enqueuer the services use, the events endpoint to mount (nil when there
// is none), and the carrier's shutdown.
//
//   - JOBS_URL unset: an in-process runner executes the jobs, holding
//     delayed ones on timers. Queued work lives in memory only.
//   - JOBS_URL set: jobs are POSTed there (a delay travels with them); the
//     queue holds and delivers them to the events endpoint.
func (app *Univents) initJobs(r *repos.Repos, notifier *database.Notifier, tx database.TxRunner) (libjobs.DelayedEnqueuer, http.Handler, func(context.Context) error, error) { //nolint:ireturn
	registry := jobs.Registry(jobs.Deps{
		Certs:            r.Certs,
		Editions:         r.Editions,
		Events:           r.Events,
		Purchases:        r.Purchases,
		Registrations:    r.Registrations,
		ProductPurchases: r.Products,
		Participations:   r.Programs,
		TicketTypes:      r.TicketTypes,
		Email:            app.emailClient,
		Notifier:         notifier,
		Tx:               tx,
	})

	if app.cfg.JobsURL != "" {
		carrier, err := httpjobs.New(app.cfg.JobsURL)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("jobs carrier: %w", err)
		}
		return carrier, httpjobs.Handler(registry), func(context.Context) error { return nil }, nil
	}

	runner := inprocess.New(registry, inprocess.Options{})
	runner.Start()
	return runner, nil, func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, jobsDrainTimeout)
		defer cancel()
		return runner.Stop(ctx)
	}, nil
}
