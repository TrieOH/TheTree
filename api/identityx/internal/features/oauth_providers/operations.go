// Package oauth_providers is the oauth providers feature: its repo, handlers
// and operations. Build wires it.
//
// The operations own everything about OAuth providers: the
// configured/enabled policy (CRUD, discovery), provider credential
// resolution, and the login flow (connect, callback, one-time login state,
// userinfo fetch, identity linking) at one interface. authn crosses it for
// the flow; handlers cross it for policy. Provider metadata and the HTTP
// client are injected at construction, so tests point both at a local
// server and never mutate a global registry.
package oauth_providers

import (
	"IdentityX/errs"
	"IdentityX/internal/authz"
	"IdentityX/internal/features/tos"
	"IdentityX/internal/tokens"
	"IdentityX/models"
	"IdentityX/ports"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"lib/crypto"
	"lib/database"
	"lib/oauth"
	"lib/telemetry"
	"log/slog"
	"time"
	"uuid"

	"golang.org/x/oauth2"
	"resty.dev/v3"
)

type Deps struct {
	Providers   ports.ProjectOAuthProvidersRepo
	LoginStates ports.OAuthLoginStatesRepo
	Projects    ports.ProjectRepo
	External    ports.ExternalIdentitiesRepo
	Actors      ports.ActorRepo
	Authz       *authz.Service
	// Tokens owns the JWT lifecycle; the callback crosses it to mint the
	// session pair for a signed-in actor.
	Tokens *tokens.Manager
	// Tos owns the terms documents and the consent ledger; connect gates
	// on it and first-time registrations pin the acceptance row through it.
	Tos *tos.Operations
	// Client and meta are the HTTP seam: the token exchange and userinfo
	// calls go through the injected resty client, and provider metadata
	// (endpoints, userinfo URL) comes from the injected registry.
	Client *resty.Client
	Meta   map[string]oauth.Provider
	// Tx opens the transactions multi-step writes run in; threaded
	// from boot, no package-level runner.
	Tx database.TxRunner
	// Vault seals the project client secrets and provider tokens at rest.
	Vault *crypto.Vault
	// PlatformCreds are IdentityX's own client credentials per provider,
	// used for logins with no project scope.
	PlatformCreds map[string]oauth.Credentials
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

// platformCreds returns IdentityX's own credentials for a provider,
// ok=false when the provider is not configured.
func (o *Operations) platformCreds(provider string) (oauth.Credentials, bool) {
	creds := o.PlatformCreds[provider]
	return creds, creds.ClientID != "" && creds.ClientSecret != ""
}

// Callback completes an OAuth login: it validates the one-time state,
// exchanges the code for provider tokens, resolves or links the identity,
// and mints the session pair through the Token-lifecycle module. The whole
// flow — state lifecycle, provider resolution, userinfo fetch, identity
// linking — lives behind this interface.
func (o *Operations) Callback(ctx context.Context, provider, code, state string) (*models.UserTokensOutput, error) {
	loginState, err := o.LoginStates.GetByState(ctx, state)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, errs.ErrOAuthStateNotFound
		}
		return nil, err
	}
	// The state is one-time-use: consume it regardless of the outcome so it
	// can never be replayed.
	err = o.LoginStates.DeleteState(ctx, loginState.ID)
	if err != nil {
		telemetry.Log().WarnContext(ctx, "failed to delete oauth login state", slog.String("state_id", loginState.ID.String()))
	}
	if string(loginState.Provider) != provider {
		return nil, errs.ErrOAuthStateMismatch
	}
	if time.Now().After(loginState.ExpiresAt) {
		return nil, errs.ErrOAuthStateExpired
	}

	resolved, err := o.resolveLoginProvider(ctx, string(loginState.Provider), loginState.ProjectID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, &errs.ProviderRemovedError{Provider: string(loginState.Provider)}
		}
		return nil, err
	}

	info, providerToken, err := o.fetchUserInfo(ctx, provider, resolved.Creds, code)
	if err != nil {
		return nil, err
	}

	encryptedAccess, err := o.Vault.EncryptPrivateKey([]byte(providerToken.AccessToken))
	if err != nil {
		return nil, err
	}
	var encryptedRefresh *string
	if providerToken.RefreshToken != "" {
		e, err := o.Vault.EncryptPrivateKey([]byte(providerToken.RefreshToken))
		if err != nil {
			return nil, err
		}
		encryptedRefresh = &e
	}
	var tokenExpiresAt *time.Time
	if !providerToken.Expiry.IsZero() {
		tokenExpiresAt = &providerToken.Expiry
	}

	identity, err := o.External.GetByProviderAndSubject(ctx, provider, info.SubString(), loginState.ProjectID)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return nil, err
	}

	var actor *models.Actor
	if identity != nil {
		actor, err = o.updateExistingIdentity(ctx, provider, info, identity, encryptedAccess, encryptedRefresh, tokenExpiresAt, loginState.ProjectID)
	} else {
		// A disabled provider allows existing identities to log back in,
		// but never new sign-ups.
		if resolved.Disabled {
			return nil, &errs.ProviderDisabledError{Provider: provider}
		}
		actor, err = o.registerNewIdentity(ctx, provider, info, encryptedAccess, encryptedRefresh, tokenExpiresAt, loginState.ProjectID)
	}
	if err != nil {
		return nil, err
	}

	return o.Tokens.Mint(ctx, actor)
}

// Connect starts an OAuth login: it resolves the scope's provider
// credentials, records a one-time login state, and returns the provider's
// authorization URL for the caller to redirect to. Provider policy lives
// in this module, so this is the flow's only contact with "is this
// provider configured/enabled".
//
// acceptedTos is the clickwrap consent flag: a project with terms only
// starts flows whose caller agreed to them, so the consent screen precedes
// the provider redirect (the callback itself has no UI to ask in).
func (o *Operations) Connect(ctx context.Context, provider string, projectID *uuid.UUID, acceptedTos bool) (string, error) {
	p, ok := o.Meta[provider]
	if !ok {
		return "", &errs.UnsupportedProviderError{Provider: provider}
	}

	if projectID != nil {
		// The project must exist before its provider row can resolve; a
		// missing project surfaces as not-found, not as a misconfigured row.
		_, err := o.Projects.GetByID(ctx, *projectID)
		if err != nil {
			return "", err
		}
		_, hasTos, err := o.Tos.CurrentTos(ctx, *projectID)
		if err != nil {
			return "", err
		}
		if hasTos && !acceptedTos {
			return "", errs.ErrTermsNotAccepted
		}
	}

	resolved, err := o.resolveLoginProvider(ctx, provider, projectID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return "", &errs.ProviderNotConfiguredError{Provider: provider, ForProject: true}
		}
		return "", err
	}

	state, err := o.createLoginState(ctx, provider, projectID)
	if err != nil {
		return "", err
	}

	return p.Config(resolved.Creds).AuthCodeURL(state), nil
}

func (o *Operations) Create(ctx context.Context, payload models.CreateOAuthProviderInput) (*models.ProjectOAuthProviders, error) {
	err := o.Authz.CheckProject(ctx, payload.ProjectID, models.ProjectRoleAdmin)
	if err != nil {
		return nil, err
	}

	encryptedSecret, err := o.Vault.EncryptPrivateKey([]byte(payload.ClientSecret))
	if err != nil {
		telemetry.Log().ErrorContext(ctx, "encrypt oauth client secret", telemetry.Err(err))
		return nil, err
	}

	created, err := o.Providers.Create(ctx, models.ProjectOAuthProviders{
		ProjectID:             payload.ProjectID,
		Provider:              payload.Provider,
		ClientID:              payload.ClientID,
		EncryptedClientSecret: encryptedSecret,
		CallbackURL:           payload.CallbackURL,
		Enabled:               true,
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (o *Operations) Delete(ctx context.Context, id uuid.UUID) error {
	row, err := o.Providers.GetByID(ctx, id)
	if err != nil {
		return err
	}
	err = o.requireProjectAdmin(ctx, row.ProjectID)
	if err != nil {
		return err
	}
	return o.Providers.Delete(ctx, id)
}

// providerOrder keeps discovery output deterministic. The registry is a map;
// the login modal renders buttons in this order.
var providerOrder = []string{"google", "github"}

// ListEnabledProviders is the public discovery route: the providers a
// project (or IdentityX itself, when projectID is nil) currently accepts
// logins with. Disabled providers are hidden so the frontend modal never
// offers a login that would be rejected.
func (o *Operations) ListEnabledProviders(ctx context.Context, projectID *uuid.UUID) ([]models.OAuthProvider, error) {
	if projectID == nil {
		var out []models.OAuthProvider
		for _, name := range providerOrder {
			if _, ok := o.platformCreds(name); ok {
				out = append(out, models.OAuthProvider(name))
			}
		}
		return out, nil
	}

	_, err := o.Projects.GetByID(ctx, *projectID)
	if err != nil {
		return nil, err
	}

	rows, err := o.Providers.ListByProject(ctx, *projectID)
	if err != nil {
		return nil, err
	}
	var out []models.OAuthProvider
	for _, row := range rows {
		if row.Enabled {
			out = append(out, row.Provider)
		}
	}
	return out, nil
}

// requireProjectAdmin gates a write on the authenticated actor holding an
// admin or owner role in the provider's project. The caller is resolved
// by the Access-check module from the context identity.
func (o *Operations) requireProjectAdmin(ctx context.Context, projectID uuid.UUID) error {
	return o.Authz.CheckProject(ctx, projectID, models.ProjectRoleAdmin)
}

func (o *Operations) ListByProject(ctx context.Context, projectID uuid.UUID) ([]models.ProjectOAuthProviders, error) {
	err := o.Authz.CheckProject(ctx, projectID, models.ProjectRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Providers.ListByProject(ctx, projectID)
}

// LoginProvider is the resolved credential scope for one OAuth provider:
// the credentials a login flow exchanges codes with, plus whether the
// provider is disabled for new sign-ups.
type LoginProvider struct {
	Provider models.OAuthProvider
	Creds    oauth.Credentials
	Disabled bool
}

// resolveLoginProvider returns the credentials and enablement for one
// provider in a scope — the platform (projectID nil, env-configured) or a
// project's configured provider row. It is the single place the login flow
// learns whether a provider is configured and enabled: connect and callback
// consult it instead of reaching into the repo, so discovery
// (ListEnabledProviders) and the flow can never drift on enablement.
//
// A missing project row surfaces as NotFound — callers decide how to phrase
// it. Project existence itself is the caller's concern.
func (o *Operations) resolveLoginProvider(ctx context.Context, provider string, projectID *uuid.UUID) (LoginProvider, error) {
	if projectID == nil {
		creds, ok := o.platformCreds(provider)
		if !ok {
			return LoginProvider{}, &errs.ProviderNotConfiguredError{Provider: provider}
		}
		return LoginProvider{Provider: models.OAuthProvider(provider), Creds: creds}, nil
	}

	row, err := o.Providers.GetByProjectAndProvider(ctx, *projectID, models.OAuthProvider(provider))
	if err != nil {
		return LoginProvider{}, err
	}

	secret, err := o.Vault.DecryptPrivateKey(row.EncryptedClientSecret)
	if err != nil {
		return LoginProvider{}, err
	}

	return LoginProvider{
		Provider: models.OAuthProvider(provider),
		Creds: oauth.Credentials{
			ClientID:     row.ClientID,
			ClientSecret: string(secret),
			RedirectURL:  row.CallbackURL,
		},
		Disabled: !row.Enabled,
	}, nil
}

func (o *Operations) SetEnabled(ctx context.Context, id uuid.UUID, enabled bool) (*models.ProjectOAuthProviders, error) {
	row, err := o.Providers.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	err = o.requireProjectAdmin(ctx, row.ProjectID)
	if err != nil {
		return nil, err
	}

	row, err = o.Providers.SetEnabled(ctx, id, enabled)
	if err != nil {
		return nil, err
	}
	return row, nil
}

func (o *Operations) Update(ctx context.Context, payload models.UpdateOAuthProviderInput) (*models.ProjectOAuthProviders, error) {
	row, err := o.Providers.GetByID(ctx, payload.ID)
	if err != nil {
		return nil, err
	}
	err = o.requireProjectAdmin(ctx, row.ProjectID)
	if err != nil {
		return nil, err
	}

	if payload.ClientID != nil {
		row, err = o.Providers.UpdateClientID(ctx, payload.ID, *payload.ClientID)
		if err != nil {
			return nil, err
		}
	}
	if payload.ClientSecret != nil {
		encryptedSecret, err := o.Vault.EncryptPrivateKey([]byte(*payload.ClientSecret))
		if err != nil {
			telemetry.Log().ErrorContext(ctx, "encrypt oauth client secret", telemetry.Err(err))
			return nil, err
		}
		row, err = o.Providers.UpdateClientSecret(ctx, payload.ID, encryptedSecret)
		if err != nil {
			return nil, err
		}
	}
	if payload.CallbackURL != nil {
		row, err = o.Providers.UpdateCallbackURL(ctx, payload.ID, *payload.CallbackURL)
		if err != nil {
			return nil, err
		}
	}

	return row, nil
}

// oauthStateTTL bounds how long a connect attempt stays valid. It must
// outlive the provider consent round trip; it is deliberately short.
const oauthStateTTL = 10 * time.Minute

// createLoginState records the connect attempt and returns the opaque state
// token the provider will hand back on the callback. Deleted there.
func (o *Operations) createLoginState(ctx context.Context, provider string, projectID *uuid.UUID) (string, error) {
	token, err := newStateToken()
	if err != nil {
		return "", err
	}
	_, err = o.LoginStates.CreateState(ctx, models.OAuthLoginState{
		State:     token,
		Provider:  models.OAuthProvider(provider),
		ProjectID: projectID,
		ExpiresAt: time.Now().Add(oauthStateTTL),
	})
	if err != nil {
		return "", err
	}
	return token, nil
}

func newStateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// updateExistingIdentity refreshes the stored provider tokens of an already
// linked identity and returns its actor. Allowed even when the provider is
// disabled — the user already has an account. The project keeps the update
// scoped: only the identity row whose actor lives in that project is touched.
func (o *Operations) updateExistingIdentity(
	ctx context.Context,
	provider string,
	info *oauth.UserInfo,
	identity *models.ActorExternalIdentities,
	encryptedAccess string,
	encryptedRefresh *string,
	tokenExpiresAt *time.Time,
	projectID *uuid.UUID,
) (*models.Actor, error) {
	_, err := o.External.UpdateTokens(ctx, models.ActorExternalIdentities{
		Provider:              models.OAuthProvider(provider),
		Subject:               info.SubString(),
		EncryptedAccessToken:  &encryptedAccess,
		EncryptedRefreshToken: encryptedRefresh,
		TokenExpiresAt:        tokenExpiresAt,
	}, projectID)
	if err != nil {
		return nil, err
	}
	actor, err := o.Actors.GetByID(ctx, identity.ActorID)
	if err != nil {
		return nil, err
	}
	// An existing identity signing in is continued use: after a version's
	// effective date the first login records the ledger row (no-op inside
	// the notice window or when already accepted). Never blocks the login.
	if actor.ProjectID != nil {
		err = o.Tos.StampUse(ctx, actor.ID, *actor.ProjectID)
		if err != nil {
			telemetry.Log().ErrorContext(ctx, "failed to stamp tos continued-use acceptance",
				slog.String("actor_id", actor.ID.String()),
				slog.String("project_id", actor.ProjectID.String()),
				telemetry.Err(err),
			)
		}
	}
	return actor, nil
}

// registerNewIdentity creates the actor and its external identity for a
// first-time OAuth login. When the flow was started from a project, the new
// actor is scoped to that project so its tokens carry the project id.
func (o *Operations) registerNewIdentity(
	ctx context.Context,
	provider string,
	info *oauth.UserInfo,
	encryptedAccess string,
	encryptedRefresh *string,
	tokenExpiresAt *time.Time,
	projectID *uuid.UUID,
) (*models.Actor, error) {
	// The account, its external identity, and its consent row commit
	// together: a compliance artifact must never exist without the ledger
	// entry that proves it. A failure rolls back everything and the
	// callback surfaces the error — no half-registered actor.
	var actor *models.Actor
	var txErr error
	err := o.Tx.WithinTx(ctx, func(ctx context.Context) error {
		actor, txErr = o.Actors.Register(ctx, models.Actor{
			ProjectID:  projectID,
			AuthMethod: models.AuthMethod(provider),
			Email:      &info.Email,
			Type:       models.HumanActorType,
			// The provider already verified the email (Google/GitHub), so the
			// account ships verified and never needs the verify link.
			VerifiedAt: new(time.Now()),
		})
		if txErr != nil {
			return txErr
		}
		_, txErr = o.External.Create(ctx, models.ActorExternalIdentities{
			ActorID:               actor.ID,
			Provider:              models.OAuthProvider(provider),
			Subject:               info.SubString(),
			Email:                 &info.Email,
			EncryptedAccessToken:  &encryptedAccess,
			EncryptedRefreshToken: encryptedRefresh,
			TokenExpiresAt:        tokenExpiresAt,
		})
		if txErr != nil {
			return txErr
		}
		if projectID == nil {
			return nil
		}
		var currentTos *models.TermsOfService
		currentTos, _, tosErr := o.Tos.CurrentTos(ctx, *projectID)
		if tosErr != nil {
			return tosErr
		}
		return o.Tos.BindRegistration(ctx, actor.ID, *projectID, currentTos)
	})
	if err != nil {
		return nil, err
	}
	return actor, nil
}

// fetchUserInfo exchanges the authorization code with the scope's
// credentials and resolves the provider's userinfo. The code exchange and
// the userinfo call both go through the injected client, so tests point it
// at a local server instead of touching the network.
func (o *Operations) fetchUserInfo(ctx context.Context, provider string, creds oauth.Credentials, code string) (*oauth.UserInfo, *oauth2.Token, error) {
	p, ok := o.Meta[provider]
	if !ok {
		return nil, nil, &errs.UnsupportedProviderError{Provider: provider}
	}

	cfg := p.Config(creds)
	exchangeCtx := context.WithValue(ctx, oauth2.HTTPClient, o.Client.Client())
	providerToken, err := cfg.Exchange(exchangeCtx, code)
	if err != nil {
		telemetry.Log().ErrorContext(ctx, "oauth code exchange failed", telemetry.Err(err))
		return nil, nil, errs.ErrOAuthCodeExchange
	}
	if providerToken == nil {
		return nil, nil, errs.ErrOAuthEmptyToken
	}

	resp, err := o.Client.R().SetContext(ctx).
		SetHeader("Authorization", "Bearer "+providerToken.AccessToken).
		Get(p.Userinfo)
	if err != nil {
		return nil, nil, err
	}
	if resp.IsStatusFailure() {
		return nil, nil, errs.ErrOAuthUserinfo
	}

	var info oauth.UserInfo
	err = json.Unmarshal(resp.Bytes(), &info)
	if err != nil {
		return nil, nil, err
	}

	if info.SubString() == "" {
		return nil, nil, errs.ErrOAuthIncompleteUserinfo
	}

	if info.Email == "" && provider == "github" {
		info.Email, err = o.fetchGitHubEmail(ctx, providerToken.AccessToken)
		if err != nil {
			return nil, nil, err
		}
	}

	if info.Email == "" {
		return nil, nil, errs.ErrOAuthIncompleteUserinfo
	}

	return &info, providerToken, nil
}

// fetchGitHubEmail resolves a GitHub user's verified primary email. GitHub
// does not include email in the userinfo payload, so it must be fetched
// separately; it uses the same injected client as the rest of the flow.
func (o *Operations) fetchGitHubEmail(ctx context.Context, accessToken string) (string, error) {
	resp, err := o.Client.R().SetContext(ctx).
		SetHeader("Authorization", "Bearer "+accessToken).
		SetHeader("Accept", "application/vnd.github+json").
		Get("https://api.github.com/user/emails")
	if err != nil {
		return "", fmt.Errorf("%w: %w", errs.ErrOAuthGitHubEmail, err)
	}
	if resp.IsStatusFailure() {
		return "", errs.ErrOAuthGitHubEmail
	}
	var emails []oauth.GitHubEmail
	err = json.Unmarshal(resp.Bytes(), &emails)
	if err != nil {
		return "", fmt.Errorf("%w: %w", errs.ErrOAuthGitHubEmail, err)
	}
	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}
	return "", fmt.Errorf("%w: no verified primary email found", errs.ErrOAuthGitHubEmail)
}
