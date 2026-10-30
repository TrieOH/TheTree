package httpserver

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"lib/authz"
)

// Authenticators verify the two credentials a request can carry. Each
// returns the context with the caller's identity on it; any error answers
// 401 with the error's message (or its mapped one).
type Authenticators struct {
	JWT    func(ctx context.Context, token string) (context.Context, error)
	APIKey func(ctx context.Context, rawKey string) (context.Context, error)
}

// The auth middlewares' own rejections: a request without the credential
// its operation needs. An authenticator's error is written as it is; the
// Backend's mapper answers both (ADR-0019).
var (
	ErrMissingBearer = errors.New("missing or invalid Authorization header")
	ErrMissingAPIKey = errors.New("missing X-API-Key header")
	ErrNoCredentials = errors.New("no authentication provided: expected X-API-Key or Authorization: Bearer <token>")
)

// Primitives builds the JWT / API key / either middlewares the spec-derived
// chains resolve against. Any prefers an API key when both are sent.
func (a Authenticators) Primitives() authz.Primitives {
	jwt := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := bearer(r)
			if !ok {
				WriteError(w, r, ErrMissingBearer)
				return
			}
			a.serve(w, r, next, a.JWT, token)
		})
	}
	apiKey := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("X-API-Key")
			if key == "" {
				WriteError(w, r, ErrMissingAPIKey)
				return
			}
			a.serve(w, r, next, a.APIKey, key)
		})
	}
	return authz.Primitives{
		JWT:    jwt,
		APIKey: apiKey,
		Any: func(next http.Handler) http.Handler {
			viaKey, viaJWT := apiKey(next), jwt(next)
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch {
				case r.Header.Get("X-API-Key") != "":
					viaKey.ServeHTTP(w, r)
				case strings.HasPrefix(r.Header.Get("Authorization"), "Bearer "):
					viaJWT.ServeHTTP(w, r)
				default:
					WriteError(w, r, ErrNoCredentials)
				}
			})
		},
	}
}

func (Authenticators) serve(w http.ResponseWriter, r *http.Request, next http.Handler,
	authenticate func(context.Context, string) (context.Context, error), credential string,
) {
	ctx, err := authenticate(r.Context(), credential)
	if err != nil {
		WriteError(w, r, err)
		return
	}
	next.ServeHTTP(w, r.WithContext(ctx))
}

func bearer(r *http.Request) (string, bool) {
	token, found := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	return token, found && token != ""
}
