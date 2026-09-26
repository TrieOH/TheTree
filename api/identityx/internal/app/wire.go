package app

import (
	"context"
	"net/http"
	"time"

	"IdentityX/internal/authz"
	"IdentityX/internal/emails"
	"IdentityX/internal/handlers"
	"IdentityX/internal/jobs"
	"IdentityX/internal/keys"
	"IdentityX/internal/repos"
	"IdentityX/internal/services"
	"IdentityX/internal/sqlc"
	"IdentityX/internal/tokens"
	libauthz "lib/authz"
	"lib/database"
	libriver "lib/river"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
)

// ── Init functions ────────────────────────────────────────────────────────

func (app *IdentityX) initRepos(q *sqlc.Queries) *repos.Repos {
	return repos.New(q)
}

func (app *IdentityX) initActionTokens(r *repos.Repos) *tokens.ActionTokenManager {
	return tokens.NewActionTokenManager(r.ActionTokens, []byte(app.cfg.HmacSecret), tokens.ActionTokenConfig{
		VerifyTTL: app.cfg.EmailVerifyTokenTTL,
		ResetTTL:  app.cfg.EmailResetTokenTTL,
	})
}

func (app *IdentityX) initTokens(r *repos.Repos) *tokens.Manager {
	return tokens.NewManager(r.CryptoKeys, r.Blacklist, r.Actors, r.Projects, tokens.Config{
		Issuer:     app.cfg.Issuer,
		AccessTTL:  app.cfg.AccessTokenLifetime,
		RefreshTTL: app.cfg.RefreshTokenLifetime,
	})
}

// initKeys constructs the Key-lifecycle module: the single owner of
// provisioning, rotation, and retirement for every scope's signing and
// encryption keys. The policy knobs are resolved here — KeyLifetime for
// the stamped expiry, RefreshTokenLifetime as the retiring grace period,
// RotateKeysJobDuration as the proactive lead window and worker interval.
func (app *IdentityX) initKeys(r *repos.Repos) *keys.Manager {
	return keys.NewManager(r.CryptoKeys, r.Projects, keys.Config{
		KeyLifetime:    app.cfg.KeyLifetime,
		RefreshTTL:     app.cfg.RefreshTokenLifetime,
		RotateInterval: app.cfg.RotateKeysJobDuration,
	})
}

func (app *IdentityX) initOperations(r *repos.Repos, tokensMgr *tokens.Manager, actionTokenMgr *tokens.ActionTokenManager, keysMgr *keys.Manager, riverClient *river.Client[pgx.Tx], tx database.TxRunner) (*services.Operations, *authz.Service) {
	authzSvc := authz.New(r.Organizations, r.Projects, r.PlatformRoles)
	sender := emails.NewSender(actionTokenMgr, app.cfg.AppURL, app.cfg.AppName, riverClient)
	tosNotifier := emails.NewTosNotifier(riverClient)
	return services.NewOperations(r, authzSvc, tokensMgr, actionTokenMgr, keysMgr, app.cfg.HmacSecret, sender, tosNotifier, tx), authzSvc
}

// initMiddlewares builds the auth primitives the spec-derived chains
// resolve against. The scope checkers come from the Access-check module:
// adding a scope is one spec x-scope line plus one entry in
// authz.ScopeCheckers. Construction stays per-backend; everything
// downstream (chain derivation, dispatch, fail-closed) is the Access-check
// and Harness modules' implementation.
func (app *IdentityX) initMiddlewares(ops *services.Operations, tokensMgr *tokens.Manager, authzSvc *authz.Service) libauthz.Primitives {
	authMW := app.SetupAuthMiddlewares(tokensMgr, ops)
	return libauthz.Primitives{
		JWT:    authMW.JWT(),
		APIKey: authMW.APIKey(),
		Any:    authMW.AnyAuth(),
		Scopes: authzSvc.ScopeCheckers(),
	}
}

func (app *IdentityX) initHandlers(ops *services.Operations) *handlers.Server {
	return handlers.NewServer(ops)
}

// initRiver migrates, starts the river client with the service's workers and
// periodic jobs, and brings up the riverui dashboard (policy lives in
// lib/river). Returns both so callers can enqueue work and mount the UI.
func (app *IdentityX) initRiver(ctx context.Context, q *sqlc.Queries, actionTokenMgr *tokens.ActionTokenManager, keysMgr *keys.Manager) (*river.Client[pgx.Tx], http.Handler, error) {
	libriver.Migrate(ctx, app.db)

	client, err := libriver.Start(ctx, app.db, libriver.NewWorkers(
		libriver.Register[jobs.CleanupBlacklistArgs](jobs.NewCleanupBlacklistWorker(q)),
		libriver.Register[jobs.CleanupActionTokensArgs](jobs.NewCleanupActionTokensWorker(actionTokenMgr)),
		libriver.Register[jobs.RotateKeysArgs](jobs.NewRotateKeysWorker(keysMgr)),
		libriver.Register[emails.SendAuthEmailArgs](jobs.NewSendAuthEmailWorker(app.emailClient, repos.NewEmailTemplates(q))),
		libriver.Register[emails.SendTosUpdateArgs](jobs.NewSendTosUpdateWorker(app.emailClient, repos.NewActors(q), repos.NewEmailTemplates(q))),
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
		// Boot provisions every scope inline (run.go), so this worker does
		// not need RunOnStart: it is the ongoing heartbeat that rotates and
		// sweeps keys while the service stays up.
		river.NewPeriodicJob(
			river.PeriodicInterval(app.cfg.RotateKeysJobDuration),
			func() (river.JobArgs, *river.InsertOpts) {
				return jobs.RotateKeysArgs{}, nil
			},
			&river.PeriodicJobOpts{RunOnStart: false},
		),
	})
	if err != nil {
		return nil, nil, err
	}

	riverUI, err := libriver.Dashboard(ctx, client)
	if err != nil {
		return nil, nil, err
	}

	return client, riverUI, nil
}
