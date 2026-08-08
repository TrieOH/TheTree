package app

import (
	"context"
	"lib/errx"
	libriver "lib/river"
	"log/slog"
	"net/http"
	"univents/internal/authz"
	"univents/internal/handlers"
	"univents/internal/repos"
	"univents/internal/services"
	certsJobs "univents/internal/services/certifications/jobs"
	"univents/internal/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"riverqueue.com/riverui"
)

// ── Wire types ────────────────────────────────────────────────────────────

type middlewares struct {
	jwt     func(http.Handler) http.Handler
	apiKey  func(http.Handler) http.Handler
	anyAuth func(http.Handler) http.Handler
}

// ── Init methods ──────────────────────────────────────────────────────────

func (app *Univents) initRepos() *repos.Repos {
	return repos.New(sqlc.New(app.db))
}

func (app *Univents) initOperations(r *repos.Repos) *services.Operations {
	authzSvc := authz.New(r.Events)
	return services.NewOperations(r, authzSvc, app.objStorage, app.idxClient, app.emailClient, app.cfg.HmacSecret, app.payssage)
}

func (app *Univents) initHandlers(ops *services.Operations) *handlers.Server {
	return handlers.NewServer(ops)
}

func (app *Univents) initMiddlewares() middlewares {
	var mw middlewares
	authMW := SetupAuthMiddlewares()

	mw.jwt = authMW.JWT()
	mw.apiKey = authMW.APIKey()
	mw.anyAuth = authMW.AnyAuth()
	return mw
}

func (app *Univents) initRiver(ctx context.Context, r *repos.Repos) (*river.Client[pgx.Tx], *riverui.Handler) {
	libriver.Migrate(ctx, app.db)

	client := libriver.NewClient(app.db, libriver.NewWorkers(
		libriver.Register(certsJobs.NewGrantCertsWorker(r.Certs, r.Editions, r.Events, app.emailClient)),
		libriver.Register(certsJobs.NewGrantCertsForOccurrenceWorker(r.Certs, r.Editions, r.Events, app.emailClient)),
	), nil, nil)
	// TODO: schedule GrantCertsForEdition on edition end and GrantCertsForOccurrence on occurrence end

	err := client.Start(ctx)
	if err != nil {
		errx.Exit(err, "failed to start river client")
	}

	riverUIHandler, err := riverui.NewHandler(&riverui.HandlerOpts{
		DevMode:   false,
		Endpoints: riverui.NewEndpoints[pgx.Tx](client, nil),
		Logger:    slog.Default(),
		Prefix:    "/riverui",
	})
	if err != nil {
		errx.Exit(err, "failed to create river ui handler")
	}
	err = riverUIHandler.Start(ctx)
	if err != nil {
		errx.Exit(err, "failed to start river ui handler")
	}

	return client, riverUIHandler
}
