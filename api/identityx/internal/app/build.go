package app

import (
	"context"
	"fmt"
	"time"

	"IdentityX/internal/authz"
	"IdentityX/internal/config"
	"IdentityX/internal/emails"
	"IdentityX/internal/features/action_tokens"
	"IdentityX/internal/features/actors"
	"IdentityX/internal/features/api_keys"
	"IdentityX/internal/features/authn"
	"IdentityX/internal/features/blacklist"
	"IdentityX/internal/features/capabilities"
	"IdentityX/internal/features/crypto_keys"
	"IdentityX/internal/features/email_templates"
	"IdentityX/internal/features/oauth_providers"
	"IdentityX/internal/features/organizations"
	"IdentityX/internal/features/platform_roles"
	"IdentityX/internal/features/profile_schemas"
	"IdentityX/internal/features/profiles"
	"IdentityX/internal/features/projects"
	"IdentityX/internal/features/tos"
	"IdentityX/internal/jobs"
	"IdentityX/internal/keys"
	"IdentityX/internal/openapi"
	"IdentityX/internal/setup"
	"IdentityX/internal/sqlc"
	"IdentityX/internal/tokens"
	"lib/backend"
	"lib/crypto"
	"lib/oauth"
	"lib/telemetry"

	"github.com/go-chi/chi/v5"
	"resty.dev/v3"
)

// Build wires IdentityX on the modules main lists (postgres, email, jobs).
// IdentityX issues the tokens, so it builds its own auth primitives rather
// than listing the identityx module. The schema is never migrated here —
// `just migrate identityx` owns it.
func Build(cfg config.Config) backend.BuildFunc {
	return func(ctx context.Context, d *backend.Deps, r chi.Router) error {
		q := sqlc.New(d.DB)
		actorsRepo := actors.NewRepo(q)
		apiKeysRepo := api_keys.NewRepo(q)
		authnRepo := authn.NewRepo(q) // also the external identities
		actionTokensRepo := action_tokens.NewRepo(q)
		blacklistRepo := blacklist.NewRepo(q)
		capabilitiesRepo := capabilities.NewRepo(q)
		cryptoKeysRepo := crypto_keys.NewRepo(q)
		emailTemplatesRepo := email_templates.NewRepo(q)
		organizationsRepo := organizations.NewRepo(q)
		oauthProvidersRepo := oauth_providers.NewRepo(q) // also the login states
		platformRolesRepo := platform_roles.NewRepo(q)
		profileSchemasRepo := profile_schemas.NewRepo(q)
		profilesRepo := profiles.NewRepo(q)
		projectsRepo := projects.NewRepo(q)
		tosRepo := tos.NewRepo(q)

		// Rotation starts one interval before a key expires, so an interval
		// as long as the lifetime would rotate every key on every run.
		if cfg.RotateKeysJobDuration >= cfg.KeyLifetime {
			return fmt.Errorf("ROTATE_KEYS_JOB_DURATION (%s) must be shorter than IDENTITY_X_KEY_LIFETIME (%s)",
				cfg.RotateKeysJobDuration, cfg.KeyLifetime)
		}

		// The master key seals every private key and OAuth secret at rest.
		vault, err := crypto.NewVault(cfg.EncryptionKey)
		if err != nil {
			return fmt.Errorf("ENCRYPTION_KEY: %w", err)
		}

		actionTokens := tokens.NewActionTokenManager(actionTokensRepo, []byte(cfg.HmacSecret), tokens.ActionTokenConfig{
			VerifyTTL: cfg.EmailVerifyTokenTTL,
			ResetTTL:  cfg.EmailResetTokenTTL,
		})

		// The Key-lifecycle module owns provisioning, rotation and
		// retirement of every scope's keys: KeyLifetime stamps the expiry,
		// RefreshTokenLifetime is the retiring grace period,
		// RotateKeysJobDuration the proactive lead window and job interval.
		keysMgr := keys.NewManager(cryptoKeysRepo, projectsRepo, keys.Config{
			KeyLifetime:    cfg.KeyLifetime,
			RefreshTTL:     cfg.RefreshTokenLifetime,
			RotateInterval: cfg.RotateKeysJobDuration,
			Vault:          vault,
		}, keys.WithOnChange(announceJWKSChange(d.Jobs)))

		// Jobs are served before the keys are checked: a key EnsureAll
		// creates or rotates at boot is announced through a job, and an
		// enqueue before Serve is dropped, which left the Worker's cached
		// JWKS without the new key.
		d.Jobs.Serve(jobs.Registry(jobs.Deps{
			Queries:      q,
			ActionTokens: actionTokens,
			Keys:         keysMgr,
			Email:        d.Email,
		}), jobs.Periodic(cfg.RotateKeysJobDuration)...)

		// Every scope has keys before the router accepts traffic; the
		// periodic RotateKeys job keeps them fresh afterwards.
		err = keysMgr.EnsureAll(ctx)
		if err != nil {
			return fmt.Errorf("ensure crypto keys: %w", err)
		}
		err = keysMgr.CheckMasterKey(ctx)
		if err != nil {
			return err
		}

		tokensMgr := tokens.NewManager(cryptoKeysRepo, blacklistRepo, actorsRepo, projectsRepo, tokens.Config{
			Issuer:     cfg.Issuer,
			AccessTTL:  cfg.AccessTokenLifetime,
			RefreshTTL: cfg.RefreshTokenLifetime,
			Vault:      vault,
		})

		// Setup is complete once any actor exists.
		setupState := setup.NewState(q.HasAnyActor)

		authzSvc := authz.New(organizationsRepo, projectsRepo, platformRolesRepo)
		sender := emails.NewSender(actionTokens, cfg.AppURL, cfg.AppName, d.Jobs)
		tosNotifier := emails.NewTosNotifier(d.Jobs)
		tosOps := tos.New(tos.Deps{
			Tos: tosRepo, Projects: projectsRepo, Actors: actorsRepo, Authz: authzSvc, Notifier: tosNotifier,
		})
		ops := operations{
			actors: actors.New(actors.Deps{Actors: actorsRepo, Projects: projectsRepo, Authz: authzSvc}),
			apiKeys: api_keys.New(api_keys.Deps{
				HmacSecret: []byte(cfg.HmacSecret), Actors: actorsRepo, APIKeys: apiKeysRepo,
				Capabilities: capabilitiesRepo, Projects: projectsRepo, Authz: authzSvc,
			}),
			authn: authn.New(authn.Deps{
				Actors: actorsRepo, Projects: projectsRepo, PlatformRoles: platformRolesRepo, Tokens: tokensMgr,
				ActionTokens: actionTokens, EmailSender: sender, Tos: tosOps, Tx: d.Tx,
			}),
			capabilities: capabilities.New(capabilities.Deps{
				Actors: actorsRepo, Capabilities: capabilitiesRepo, Projects: projectsRepo, Authz: authzSvc,
			}),
			emailTemplates: email_templates.New(email_templates.Deps{Templates: emailTemplatesRepo, Authz: authzSvc}),
			oauthProviders: oauth_providers.New(oauth_providers.Deps{
				Providers: oauthProvidersRepo, LoginStates: oauthProvidersRepo, Projects: projectsRepo,
				External: authnRepo, Actors: actorsRepo, Authz: authzSvc, Tokens: tokensMgr, Tos: tosOps,
				Client: resty.New().SetTimeout(15 * time.Second), Tx: d.Tx, Vault: vault,
				Meta: oauth.NewRegistry(map[string]string{
					"google": cfg.GoogleRedirectURI,
					"github": cfg.GithubRedirectURI,
				}),
				PlatformCreds: map[string]oauth.Credentials{
					"google": {ClientID: cfg.GoogleClientID, ClientSecret: cfg.GoogleClientSecret},
					"github": {ClientID: cfg.GithubClientID, ClientSecret: cfg.GithubClientSecret},
				},
			}),
			organizations: organizations.New(organizations.Deps{
				Projects: projectsRepo, Actors: actorsRepo, Orgs: organizationsRepo, Keys: keysMgr, Authz: authzSvc, Tx: d.Tx,
			}),
			profileSchemas: profile_schemas.New(profile_schemas.Deps{Schemas: profileSchemasRepo, Projects: projectsRepo, Authz: authzSvc}),
			profiles: profiles.New(profiles.Deps{
				Profiles: profilesRepo, Schemas: profileSchemasRepo, Actors: actorsRepo, Authz: authzSvc,
			}),
			projects: projects.New(projects.Deps{
				Projects: projectsRepo, Actors: actorsRepo, Keys: keysMgr, Authz: authzSvc, Tx: d.Tx,
			}),
			tos: tosOps,
		}

		chains, err := d.Chains(authPrimitives(tokensMgr, ops.apiKeys, authzSvc))
		if err != nil {
			return err
		}
		err = guardSetup(chains, setupState)
		if err != nil {
			return err
		}
		openapi.Mount(r, newServer(ops, setupState), chains)
		return nil
	}
}

// announceJWKSChange tells whoever caches the JWKS that the keys changed: on
// Cloudflare the Worker in front of the container drops its copy. In-process
// nothing caches it. Build serves the jobs before any key can change, so a
// failed enqueue is a real loss: the Worker can keep the old key list for up
// to a week, and other APIs refuse tokens signed with the new key.
func announceJWKSChange(j backend.Jobs) func(context.Context) {
	return func(ctx context.Context) {
		err := j.Enqueue(ctx, jobs.JWKSChangedArgs{})
		if err != nil {
			telemetry.Log().ErrorContext(ctx, "could not announce a JWKS change; the Worker's cached JWKS is stale", telemetry.Err(err))
		}
	}
}
