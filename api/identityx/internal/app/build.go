package app

import (
	"context"
	"fmt"

	"IdentityX/internal/authz"
	"IdentityX/internal/config"
	"IdentityX/internal/emails"
	"IdentityX/internal/handlers"
	"IdentityX/internal/jobs"
	"IdentityX/internal/keys"
	"IdentityX/internal/repos"
	"IdentityX/internal/services"
	"IdentityX/internal/setup"
	"IdentityX/internal/sqlc"
	"IdentityX/internal/tokens"
	"lib/backend"

	"github.com/go-chi/chi/v5"
)

// Build wires IdentityX on the modules main lists (postgres, email, jobs).
// IdentityX issues the tokens, so it builds its own auth primitives rather
// than listing the identityx module. The schema is never migrated here —
// `just identityx-goose` owns it.
func Build(cfg config.Config) backend.BuildFunc {
	return func(ctx context.Context, d *backend.Deps, r chi.Router) error {
		q := sqlc.New(d.DB)
		repos := repos.New(q)

		actionTokens := tokens.NewActionTokenManager(repos.ActionTokens, []byte(cfg.HmacSecret), tokens.ActionTokenConfig{
			VerifyTTL: cfg.EmailVerifyTokenTTL,
			ResetTTL:  cfg.EmailResetTokenTTL,
		})

		// The Key-lifecycle module owns provisioning, rotation and
		// retirement of every scope's keys: KeyLifetime stamps the expiry,
		// RefreshTokenLifetime is the retiring grace period,
		// RotateKeysJobDuration the proactive lead window and job interval.
		keysMgr := keys.NewManager(repos.CryptoKeys, repos.Projects, keys.Config{
			KeyLifetime:    cfg.KeyLifetime,
			RefreshTTL:     cfg.RefreshTokenLifetime,
			RotateInterval: cfg.RotateKeysJobDuration,
		})
		// Every scope has keys before the router accepts traffic; the
		// periodic RotateKeys job keeps them fresh afterwards.
		err := keysMgr.EnsureAll(ctx)
		if err != nil {
			return fmt.Errorf("ensure crypto keys: %w", err)
		}

		tokensMgr := tokens.NewManager(repos.CryptoKeys, repos.Blacklist, repos.Actors, repos.Projects, tokens.Config{
			Issuer:     cfg.Issuer,
			AccessTTL:  cfg.AccessTokenLifetime,
			RefreshTTL: cfg.RefreshTokenLifetime,
		})

		// Setup is complete once any actor exists.
		setupState := setup.NewState(q.HasAnyActor)

		d.Jobs.Serve(jobs.Registry(jobs.Deps{
			Queries:      q,
			ActionTokens: actionTokens,
			Keys:         keysMgr,
			Email:        d.Email,
		}), jobs.Periodic(cfg.RotateKeysJobDuration)...)

		authzSvc := authz.New(repos.Organizations, repos.Projects, repos.PlatformRoles)
		sender := emails.NewSender(actionTokens, cfg.AppURL, cfg.AppName, d.Jobs)
		tosNotifier := emails.NewTosNotifier(d.Jobs)
		ops := services.NewOperations(repos, authzSvc, tokensMgr, actionTokens, keysMgr, cfg.HmacSecret, sender, tosNotifier, d.Tx)

		chains, err := d.ChainsWith(authPrimitives(tokensMgr, ops, authzSvc), chainOptions(setupState))
		if err != nil {
			return err
		}
		mountStrict(r, handlers.NewServer(ops, setupState), chains)
		return nil
	}
}
