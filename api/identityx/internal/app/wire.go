package app

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"IdentityX/internal/authz"
	"IdentityX/internal/emails"
	"IdentityX/internal/handlers"
	"IdentityX/internal/jobs"
	"IdentityX/internal/repos"
	"IdentityX/internal/services"
	"IdentityX/internal/sqlc"
	"IdentityX/internal/tokens"
	"lib/errx"
	libriver "lib/river"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"riverqueue.com/riverui"
)

// ── Wire types ────────────────────────────────────────────────────────────

type middlewares struct {
	jwtAuth    func(http.Handler) http.Handler
	apiKeyAuth func(http.Handler) http.Handler
	anyAuth    func(http.Handler) http.Handler
}

// ── Init functions ────────────────────────────────────────────────────────

func (app *IdentityX) initRepos(q *sqlc.Queries) *repos.Repos {
	return repos.New(q)
}

func (app *IdentityX) initTokens(r *repos.Repos) *tokens.Manager {
	return tokens.NewManager(r.CryptoKeys, r.Blacklist, r.Actors, r.Projects, tokens.Config{
		Issuer:     app.cfg.Issuer,
		AccessTTL:  app.cfg.AccessTokenLifetime,
		RefreshTTL: app.cfg.RefreshTokenLifetime,
	})
}

func (app *IdentityX) initOperations(r *repos.Repos, tokensMgr *tokens.Manager, riverClient *river.Client[pgx.Tx]) *services.Operations {
	authzSvc := authz.New(r.Organizations, r.Projects, r.PlatformRoles)
	sender := emails.NewSender(
		r.ActionTokens,
		[]byte(app.cfg.HmacSecret),
		app.cfg.EmailVerifyTokenTTL,
		app.cfg.EmailResetTokenTTL,
		app.cfg.AppURL,
		app.cfg.AppName,
		riverClient,
	)
	return services.NewOperations(r, authzSvc, tokensMgr, app.cfg.HmacSecret, sender)
}

func (app *IdentityX) initMiddlewares(r *repos.Repos, tokensMgr *tokens.Manager) middlewares {
	var mw middlewares
	authMW := app.SetupAuthMiddlewares(tokensMgr, r.APIKeys, r.Actors, r.Capabilities)
	mw.jwtAuth = authMW.JWT()
	mw.apiKeyAuth = authMW.APIKey()
	mw.anyAuth = authMW.AnyAuth()
	return mw
}

func (app *IdentityX) initHandlers(ops *services.Operations) *handlers.Server {
	return handlers.NewServer(ops)
}

// initRiver migrates, starts the river client with the service's workers and
// periodic jobs, and brings up the riverui dashboard. Returns both so callers
// can enqueue work and mount the UI handler.
func (app *IdentityX) initRiver(ctx context.Context, q *sqlc.Queries) (*river.Client[pgx.Tx], *riverui.Handler) {
	libriver.Migrate(ctx, app.db)

	client := libriver.NewClient(app.db, libriver.NewWorkers(
		libriver.Register[jobs.CreateCryptoKeyArgs](jobs.NewCreateCryptoKeyWorker(q)),
		libriver.Register[jobs.CleanupBlacklistArgs](jobs.NewCleanupBlacklistWorker(q)),
		libriver.Register[jobs.CleanupActionTokensArgs](jobs.NewCleanupActionTokensWorker(q)),
		libriver.Register[emails.SendAuthEmailArgs](jobs.NewSendAuthEmailWorker(app.emailClient, repos.NewEmailTemplates(q))),
	), nil, []*river.PeriodicJob{
		river.NewPeriodicJob(
			river.PeriodicInterval(5*time.Minute),
			func() (river.JobArgs, *river.InsertOpts) {
				return jobs.CleanupBlacklistArgs{}, nil
			},
			&river.PeriodicJobOpts{RunOnStart: true},
		),
		river.NewPeriodicJob(
			river.PeriodicInterval(5*time.Minute),
			func() (river.JobArgs, *river.InsertOpts) {
				return jobs.CleanupActionTokensArgs{}, nil
			},
			&river.PeriodicJobOpts{RunOnStart: true},
		),
	})

	err := client.Start(ctx)
	if err != nil {
		errx.Exit(err, "failed to start river client")
	}

	riverUIHandler, err := riverui.NewHandler(&riverui.HandlerOpts{
		DevMode:                  false,
		Endpoints:                riverui.NewEndpoints[pgx.Tx](client, nil),
		Logger:                   slog.Default(),
		Prefix:                   "/riverui",
		JobListHideArgsByDefault: true,
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
