package app

import (
	"context"
	"lib/api_keys"
	"net/http"
	"time"

	"IdentityX/internal/database/sqlc"
	"IdentityX/internal/jobs"
	"IdentityX/models"
	"lib/crypto"
	"lib/database"
	"lib/errx"
	"lib/validator"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
	"github.com/MintzyG/fun/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
)

func SetupFUN() {
	fun.SetConfig(fun.Config{
		MaxTraceSize:         50,
		ResponseSizeLimit:    10 * 1024 * 1024,
		MaxInterceptorAmount: 20,
		DefaultContentType:   "application/json",
		EnableSizeValidation: true,
		DefaultModule:        "IdentityX-API",
	})

	v := validator.SetupValidator()
	bind.SetValidator(v)
	fun.SetPathParamFunc(func(r *http.Request, key string) string {
		return chi.URLParam(r, key)
	})
}

func SetupConstraintMessages() {
	database.SetConstraintErrorRegistry(database.ConstraintRegistry{
		// actors
		"chk_actors_type":                      "actor type must be one of: human, service, machine",
		"chk_actors_auth_method":               "auth method must be one of: api_key, password, google, github",
		"chk_actors_email_required_for_humans": "email is required for human actors",
		"uniq_email_per_scope_per_method":      "an account with this email already exists for this scope and auth method",

		// actor_external_identities
		"chk_actor_external_identities_provider": "external identity provider must be one of: google, github",
		"uniq_external_identity":                 "this external identity is already linked to another account",

		// organizations
		"uniq_organizations_slug": "an organization with this slug already exists",

		// org_members
		"chk_org_members_role": "organization role must be one of: owner, admin, member",

		// projects
		"uniq_projects_slug":   "a project with this slug already exists",
		"uniq_verified_domain": "this domain is already verified by another project",

		// project_members
		"chk_project_members_role": "project role must be one of: owner, admin, member",

		// project_oauth_providers
		"chk_project_oauth_providers_provider": "OAuth provider must be one of: google, github",
		"uniq_project_oauth_provider":          "this OAuth provider is already configured for this project",

		// platform_roles
		"chk_platform_roles_role": "platform role must be one of: super_admin, admin, support",

		// api_keys
		"uniq_idx_api_keys_key_prefix": "an API key with this prefix already exists",

		// crypto_keys
		"chk_crypto_keys_type":   "crypto key type must be one of: encryption, signing",
		"chk_crypto_keys_status": "crypto key status must be one of: active, retiring, retired, revoked",

		// blacklist_entries
		"chk_blacklist_entries_type":         "blacklist entry type must be one of: actor, token, api_key, email, ip",
		"uniq_blacklist_target_type_project": "this target is already blacklisted for this scope",
	})
}

func EnsureKeysExist(ctx context.Context, db *pgxpool.Pool, riverClient *river.Client[pgx.Tx]) {
	q := sqlc.New(db)
	for _, keyType := range []string{"signing", "encryption"} {
		exists, err := q.HasActiveCryptoKey(ctx, sqlc.HasActiveCryptoKeyParams{Type: keyType})
		if err != nil {
			errx.Exit(err, "failed to check global "+keyType+" key")
		}
		if !exists {
			if _, err = riverClient.Insert(ctx, jobs.CreateCryptoKeyArgs{KeyType: keyType}, nil); err != nil {
				errx.Exit(err, "failed to enqueue global "+keyType+" key creation")
			}
		}
	}

	projects, err := q.ListProjects(ctx)
	if err != nil {
		errx.Exit(err, "failed to list projects for key check")
	}

	for _, p := range projects {
		pid := p.ID
		for _, keyType := range []string{"signing", "encryption"} {
			exists, err := q.HasActiveCryptoKey(ctx, sqlc.HasActiveCryptoKeyParams{
				ProjectID: &pid,
				Type:      keyType,
			})
			if err != nil {
				errx.Exit(err, "failed to check "+keyType+" key for project "+pid.String())
			}
			if !exists {
				if _, err = riverClient.Insert(ctx, jobs.CreateCryptoKeyArgs{
					ProjectID: &pid,
					KeyType:   keyType,
				}, nil); err != nil {
					errx.Exit(err, "failed to enqueue "+keyType+" key creation for project "+pid.String())
				}
			}
		}
	}
}

func SetupAuthMiddlewares(rt runtime) *middlewares.Middleware[*models.AccessClaims] {
	keyFunc := func(ctx context.Context, tokenStr string) (*models.AccessClaims, error) {
		claims := &models.AccessClaims{}
		token, err := crypto.OpenUnverified(tokenStr, claims)
		if err != nil {
			return nil, err
		}
		kid, ok := token.Header["kid"].(string)
		if !ok || kid == "" {
			return nil, fun.ErrUnauthorized("missing kid")
		}
		keyID, err := uuid.Parse(kid)
		if err != nil {
			return nil, fun.ErrUnauthorized("invalid kid")
		}
		cryptoKey, err := rt.repos.cryptoKeys.GetByID(ctx, keyID)
		if fun.Is(err, fun.CodeNotFound) {
			return nil, fun.ErrUnauthorized("outdated token")
		}
		if err != nil {
			return nil, err
		}
		if cryptoKey.Status == "revoked" {
			return nil, fun.ErrUnauthorized("token signing key revoked")
		}
		_, err = crypto.VerifyToken(tokenStr, cryptoKey.PublicKey, claims)
		if err != nil {
			return nil, fun.ErrUnauthorized("invalid access token")
		}
		return claims, nil
	}

	jwtHook := func(ctx context.Context, claims *models.AccessClaims) (context.Context, error) {
		identity := &models.Identity{
			Sub: models.SubjectFromAccessSub(claims.Sub),
			Cred: models.Credential{
				Type: models.TokenCredentialType,
			},
		}
		return models.WithIdentity(ctx, identity), nil
	}

	apiKeyHook := func(ctx context.Context, rawKey string) (context.Context, error) {
		if len(rawKey) < api_keys.ApiKeyPrefixLen {
			return nil, fun.ErrUnauthorized("invalid api key")
		}
		body, err := api_keys.StripKeyHeader(rawKey)
		if err != nil || len(body) < api_keys.ApiKeyPrefixLen {
			return nil, fun.ErrUnauthorized("invalid api key")
		}
		prefix := body[:api_keys.ApiKeyPrefixLen]

		apiKey, err := rt.repos.apiKeys.GetByPrefix(ctx, prefix)
		if err != nil {
			return nil, fun.ErrUnauthorized("invalid api key")
		}

		if !api_keys.VerifyAPIKey(rawKey, apiKey.KeyHash) {
			return nil, fun.ErrUnauthorized("invalid api key")
		}

		if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
			return nil, fun.ErrUnauthorized("api key expired")
		}

		actor, err := rt.repos.actors.GetByID(ctx, apiKey.ActorID)
		if err != nil {
			return nil, fun.ErrUnauthorized("invalid api key")
		}

		ctx = models.WithIdentity(ctx, &models.Identity{
			Sub: models.Subject{
				ID:           actor.ID,
				ProjectID:    actor.ProjectID,
				Email:        actor.Email,
				Type:         actor.Type,
				Capabilities: nil,
				Metadata:     nil,
			},
			Cred: models.Credential{
				ID:   &apiKey.ID,
				Type: models.ApiKeyCredentialType,
				Raw:  rawKey,
			},
		})
		return ctx, nil
	}
	return middlewares.New[*models.AccessClaims](keyFunc, jwtHook, apiKeyHook)
}

// ClientOnly returns a middleware that rejects unauthenticated requests
// and requests from project-scoped actors. It checks that a valid identity
// exists in the context and that the identity's subject carries a nil
// ProjectID — meaning the actor is an IdentityX platform-level client
// (human, service, or machine) rather than a project-scoped one.
//
// Use it after an auth middleware to guard routes that should only be
// reachable by platform-level IdentityX clients.
func ClientOnly() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ident, err := models.RequireIdentity(r.Context())
			if fun.Bail(w, err) {
				return
			}
			if ident.Sub.ProjectID != nil {
				fun.Unauthorized("platform-level authentication required").Send(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// ProjectClientOnly returns a middleware that rejects requests that are
// not both authenticated AND scoped to a specific project. It checks
// that a valid identity exists in the context and that the identity's
// subject carries a non-nil ProjectID — meaning the actor is acting
// within a project context (e.g. workspace API keys, project service
// accounts).
//
// Use it after an auth middleware to guard routes that should only be
// reachable by project-scoped actors.
func ProjectClientOnly() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ident, err := models.RequireIdentity(r.Context())
			if fun.Bail(w, err) {
				return
			}
			if ident.Sub.ProjectID == nil {
				fun.Unauthorized("project-scoped authentication required").Send(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
