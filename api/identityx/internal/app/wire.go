package app

import (
	"context"
	"fmt"
	"time"

	"IdentityX/internal/authz"
	"IdentityX/internal/config"
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
	"lib/email"
	"lib/email/ses"
	libjobs "lib/jobs"
	"lib/jobs/inprocess"
	"lib/jobs/sqsjobs"
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
// RotateKeysJobDuration as the proactive lead window and job interval.
func (app *IdentityX) initKeys(r *repos.Repos) *keys.Manager {
	return keys.NewManager(r.CryptoKeys, r.Projects, keys.Config{
		KeyLifetime:    app.cfg.KeyLifetime,
		RefreshTTL:     app.cfg.RefreshTokenLifetime,
		RotateInterval: app.cfg.RotateKeysJobDuration,
	})
}

func (app *IdentityX) initOperations(r *repos.Repos, tokensMgr *tokens.Manager, actionTokenMgr *tokens.ActionTokenManager, keysMgr *keys.Manager, enqueuer libjobs.Enqueuer, tx database.TxRunner) (*services.Operations, *authz.Service) {
	authzSvc := authz.New(r.Organizations, r.Projects, r.PlatformRoles)
	sender := emails.NewSender(actionTokenMgr, app.cfg.AppURL, app.cfg.AppName, enqueuer)
	tosNotifier := emails.NewTosNotifier(enqueuer)
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

// initEmail builds the configured email transport.
func (app *IdentityX) initEmail(ctx context.Context) (email.Sender, error) { //nolint:ireturn
	switch app.cfg.EmailTransport {
	case config.SESTransport:
		sender, err := ses.New(ctx, app.cfg.EmailFrom)
		if err != nil {
			return nil, fmt.Errorf("email transport: %w", err)
		}
		return sender, nil
	default:
		return email.NewClient(app.cfg.ToEmailConfig()), nil
	}
}

// jobsDrainTimeout bounds how long hosted shutdown waits for queued jobs.
const jobsDrainTimeout = 15 * time.Second

// initJobs picks the background-job carrier for the runtime mode and
// returns the enqueuer the services use plus its shutdown.
//
//   - hosted: an in-process runner executes the jobs and runs the periodic
//     ones on tickers. Queued work lives in memory only.
//   - managed: jobs go to the SQS queue; the worker role drains it through
//     the /events endpoint, and the scheduler fires the periodic ones there.
func (app *IdentityX) initJobs(ctx context.Context, registry *libjobs.Registry) (libjobs.Enqueuer, func(context.Context) error, error) { //nolint:ireturn
	if app.cfg.Managed() {
		enqueuer, err := sqsjobs.New(ctx, app.cfg.JobsQueueURL)
		if err != nil {
			return nil, nil, fmt.Errorf("jobs carrier: %w", err)
		}
		return enqueuer, func(context.Context) error { return nil }, nil
	}

	runner := inprocess.New(registry, inprocess.Options{})
	runner.Start(jobs.Periodic(app.cfg.RotateKeysJobDuration)...)
	return runner, func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, jobsDrainTimeout)
		defer cancel()
		return runner.Stop(ctx)
	}, nil
}
