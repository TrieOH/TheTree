// Package identityx is the kit module for Backends that consume IdentityX:
// the SDK client, and the one implementation of the auth primitives every
// consumer authenticates requests with (JWT, API key, either). It replaces
// the per-Backend auth.go copies, which had drifted — Payssage's hardcoded
// the API-key introspect host.
//
// IdentityX itself issues the tokens, so it does not list this module: it
// builds its own primitives and passes them to Deps.Chains.
package identityx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"lib/authz"
	"lib/backend"
	libconfig "lib/config"

	idx "sdk/identityx"

	fm "github.com/MintzyG/fun/middlewares"
)

// introspectTimeout bounds one API-key introspection call.
const introspectTimeout = 15 * time.Second

// The auth hooks' rejections. fun's auth middleware answers any hook error
// with a 401 carrying its message, so these never reach the error map.
var (
	ErrEmailNotVerified   = errors.New("verified email required")
	ErrAPIKeysNotAccepted = errors.New("this service does not accept api keys")
	ErrInvalidAPIKey      = errors.New("invalid api key")
)

// Option tunes the module per Backend.
type Option func(*options)

type options struct {
	debug                bool
	requireVerifiedEmail bool
	acceptAPIKeys        bool
	httpClient           *http.Client
}

// Debug turns on the SDK client's request logging.
func Debug(on bool) Option {
	return func(o *options) { o.debug = on }
}

// RequireVerifiedEmail rejects JWT-authenticated callers whose email is not
// verified (401, ErrEmailNotVerified). Public routes are unaffected.
func RequireVerifiedEmail(on bool) Option {
	return func(o *options) { o.requireVerifiedEmail = on }
}

// AcceptAPIKeys authenticates X-API-KEY callers by introspecting the key
// against IdentityX (at the configured IDENTITY_X_URL). Without it, an API
// key is refused (401, ErrAPIKeysNotAccepted): the Backend has no public API.
func AcceptAPIKeys() Option {
	return func(o *options) { o.acceptAPIKeys = true }
}

// HTTPClient overrides the client used for introspection. Tests only.
func HTTPClient(c *http.Client) Option {
	return func(o *options) { o.httpClient = c }
}

// Module bootstraps the IdentityX SDK client and fills Deps.IDX with it and
// the auth primitives.
func Module(cfg libconfig.IdentityX, opts ...Option) backend.Module {
	o := options{httpClient: &http.Client{Timeout: introspectTimeout}}
	for _, opt := range opts {
		opt(&o)
	}
	return backend.Module{
		Name: "identityx",
		Open: func(ctx context.Context, d *backend.Deps) (backend.Close, error) {
			client, err := idx.Bootstrap(ctx, libconfig.IdentityXConfig(cfg, o.debug))
			if err != nil {
				return nil, err
			}
			d.IDX = &backend.IdentityX{
				Client: client,
				Auth:   primitives(client, cfg.URL, o),
			}
			return backend.NoClose, nil
		},
	}
}

// primitives builds the JWT / API key / either middlewares.
func primitives(client *idx.Client, baseURL string, o options) authz.Primitives {
	keyFunc := func(ctx context.Context, tokenStr string) (*idx.AccessClaims, error) {
		return client.Tokens.VerifyAccessToken(ctx, tokenStr)
	}
	mw := fm.New[*idx.AccessClaims](keyFunc, jwtHook(o.requireVerifiedEmail), apiKeyHook(baseURL, o))
	return authz.Primitives{
		JWT:    mw.JWT(),
		APIKey: mw.APIKey(),
		Any:    mw.AnyAuth(),
	}
}

// jwtHook writes the token's subject into the context as the caller's
// identity.
func jwtHook(requireVerified bool) fm.JWTHook[*idx.AccessClaims] {
	return func(ctx context.Context, claims *idx.AccessClaims) (context.Context, error) {
		if requireVerified && claims.Sub.VerifiedAt == nil {
			return ctx, ErrEmailNotVerified
		}
		return idx.WithIdentity(ctx, &idx.Identity{
			Sub:  idx.SubjectFromAccessSub(claims.Sub),
			Cred: idx.Credential{Type: "token"},
		}), nil
	}
}

// apiKeyHook refuses API keys, or introspects them when the Backend opted
// in with AcceptAPIKeys.
func apiKeyHook(baseURL string, o options) fm.APIKeyHook {
	if !o.acceptAPIKeys {
		return func(ctx context.Context, _ string) (context.Context, error) {
			return ctx, ErrAPIKeysNotAccepted
		}
	}
	introspectURL := strings.TrimSuffix(baseURL, "/") + "/auth/introspect"
	return func(ctx context.Context, rawKey string) (context.Context, error) {
		ident, err := introspect(ctx, o.httpClient, introspectURL, rawKey)
		if err != nil {
			return ctx, err
		}
		return idx.WithIdentity(ctx, ident), nil
	}
}

// introspect asks IdentityX who owns rawKey. IdentityX answers with the
// identity. A non-2xx answer is ErrInvalidAPIKey; a transport failure is
// returned as is.
func introspect(ctx context.Context, client *http.Client, url, rawKey string) (*idx.Identity, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("introspect: build request: %w", err)
	}
	req.Header.Set("X-API-KEY", rawKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("introspect: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode/100 != 2 {
		return nil, ErrInvalidAPIKey
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("introspect: read: %w", err)
	}
	var identity idx.Identity
	err = json.Unmarshal(body, &identity)
	if err != nil {
		return nil, fmt.Errorf("introspect: decode: %w", err)
	}
	return &identity, nil
}
