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
	"payssage/internal/authz"
	"payssage/internal/features/providers"
	"payssage/internal/handlers"
	"payssage/internal/jobs"
	providers2 "payssage/internal/providers"
	"payssage/internal/repos"
	"payssage/internal/services"
	"payssage/internal/sqlc"
	"payssage/ports"
)

// ── Init functions ────────────────────────────────────────────────────────

func (app *Payssage) initRepos(q *sqlc.Queries) *repos.Repos {
	return repos.New(q)
}

func (app *Payssage) initProviders(r *repos.Repos) {
	mercadoPago := providers.NewMercadoPago(app.cfg.MercadoPagoConfig, r.Intents, r.Collectors, r.Sellers, r.Wallets, app.httpClient)

	providers2.PayssageProviders.OAuth = map[providers2.AvailableProviders]ports.OAuthProvider{
		providers2.MercadoPagoProvider: mercadoPago,
	}

	providers2.PayssageProviders.Payments = map[providers2.AvailableProviders]ports.PaymentAbstractionLayer{
		providers2.MercadoPagoProvider: mercadoPago,
	}

	providers2.PayssageProviders.Webhooks = map[providers2.AvailableProviders]ports.WebhookProvider{
		providers2.MercadoPagoProvider: mercadoPago,
	}
}

func (app *Payssage) initOperations(enqueuer libjobs.Enqueuer, r *repos.Repos, tx database.TxRunner) *services.Operations {
	authzSvc := authz.New(r.Organizations, r.Wallets)
	return services.NewOperations(r, authzSvc, enqueuer, app.idxClient, tx)
}

// jobsDrainTimeout bounds how long shutdown waits for queued jobs.
const jobsDrainTimeout = 15 * time.Second

// lateEnqueuer lets the job handlers hold the enqueuer before the carrier
// exists: the in-process runner is built from the registry, and the redrive
// handler in that registry enqueues through the runner.
type lateEnqueuer struct {
	libjobs.Enqueuer
}

// initJobs builds the job registry, picks the carrier, and returns the
// enqueuer the services use, the events endpoint to mount (nil when there
// is none), and the carrier's shutdown.
//
//   - JOBS_URL unset: an in-process runner executes the jobs and runs the
//     periodic ones on tickers. Queued work lives in memory only.
//   - JOBS_URL set: jobs are POSTed there; the queue delivers them, and the
//     scheduler fires the periodic ones, at the events endpoint.
func (app *Payssage) initJobs(r *repos.Repos) (libjobs.Enqueuer, http.Handler, func(context.Context) error, error) { //nolint:ireturn
	enqueuer := &lateEnqueuer{}
	registry := jobs.Registry(jobs.Deps{
		Deliveries: r.WebhookDeliveries,
		Events:     r.WebhookEvents,
		Endpoints:  r.WebhookEndpoints,
		HTTPClient: app.httpClient,
		Enqueuer:   enqueuer,
	})

	if app.cfg.JobsURL != "" {
		carrier, err := httpjobs.New(app.cfg.JobsURL)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("jobs carrier: %w", err)
		}
		enqueuer.Enqueuer = carrier
		return carrier, httpjobs.Handler(registry), func(context.Context) error { return nil }, nil
	}

	runner := inprocess.New(registry, inprocess.Options{})
	enqueuer.Enqueuer = runner
	runner.Start(jobs.Periodic()...)
	return runner, nil, func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, jobsDrainTimeout)
		defer cancel()
		return runner.Stop(ctx)
	}, nil
}

// initMiddlewares builds the auth primitives the spec-derived chains
// resolve against; construction stays per-backend, everything downstream
// (chain derivation, dispatch, fail-closed) is the Access-check and
// Harness modules' implementation.
func (app *Payssage) initMiddlewares() libauthz.Primitives {
	authMW := app.setupAuthMiddlewares()
	return libauthz.Primitives{
		JWT:    authMW.JWT(),
		APIKey: authMW.APIKey(),
		Any:    authMW.AnyAuth(),
	}
}

func (app *Payssage) initHandlers(ops *services.Operations) *handlers.Server {
	return handlers.NewServer(ops)
}
