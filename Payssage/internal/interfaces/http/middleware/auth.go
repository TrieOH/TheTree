package middleware

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/authz"
	"TriePayments/internal/shared/errx"
	"errors"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/TrieOH/goauth-sdk-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
)

type AuthMiddleware struct {
	gaClient   goauth.Client
	apiKeys    domain.ApiKeysRepo
	workspaces domain.WorkspaceRepo
	tracer     trace.Tracer
}

func NewAuthMiddleware(
	gaClient *goauth.Client,
	apiKeys domain.ApiKeysRepo,
	workspaces domain.WorkspaceRepo,
	tracer trace.Tracer,
) *AuthMiddleware {
	return &AuthMiddleware{
		gaClient:   *gaClient,
		apiKeys:    apiKeys,
		workspaces: workspaces,
		tracer:     tracer,
	}
}

// Auth is a middleware function that checks for valid access and refresh tokens.
// It validates the tokens, checks if the refresh token is revoked, and creates a principal from the tokens.
// The principal is then added to the request context.
func (mw *AuthMiddleware) Auth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx, span := mw.tracer.Start(ctx, "Middleware.Auth")
			defer span.End()

			var rs = resp.BadRequest().WithCode(400).WithModule("AuthMW")
			var err error
			defer func() {
				span.SetAttributes(attribute.Bool("success", err == nil))
			}()

			var accessTokenCookie *http.Cookie
			accessTokenCookie, err = r.Cookie("access_token")
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					rs.WithCode(401).WithMsg(errx.NotFound("access_token cookie").Error()).Send(w)
					return
				}
				rs.WithCode(401).WithMsg(errx.Invalid("access_token").Error()).Send(w)
				return
			}

			accessToken := accessTokenCookie.Value
			token, err := mw.gaClient.Tokens.ValidateToken(ctx, accessToken)
			if err != nil {
				resp.Unauthorized(err.Error()).WithModule("AuthMW").Send(w)
				return
			}

			subject, err := authz.GetSubjectFromToken(token)
			if err != nil {
				resp.Unauthorized(err.Error()).WithModule("AuthMW").Send(w)
				return
			}

			ctx = authz.WithSubject(ctx, subject)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (mw *AuthMiddleware) APIKey() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx, span := mw.tracer.Start(ctx, "Middleware.APIKey")
			defer span.End()

			var err error
			defer func() {
				span.SetAttributes(attribute.Bool("success", err == nil))
			}()

			rawKey := r.Header.Get("X-API-Key")
			if rawKey == "" {
				resp.Unauthorized("missing X-API-Key header").WithModule("AuthMW").Send(w)
				return
			}

			// extract prefix (first 11 chars: "tp_" + 8 hex)
			if len(rawKey) < 11 {
				resp.Unauthorized("invalid api key format").WithModule("AuthMW").Send(w)
				return
			}
			prefix := rawKey[:11]

			candidates, err := mw.apiKeys.GetByPrefix(ctx, prefix)
			if err != nil || len(candidates) == 0 {
				resp.Unauthorized("invalid api key").WithModule("AuthMW").Send(w)
				return
			}

			var matched *domain.APIKey
			for _, candidate := range candidates {
				if err := bcrypt.CompareHashAndPassword([]byte(candidate.KeyHash), []byte(rawKey)); err == nil {
					matched = &candidate
					break
				}
			}
			if matched == nil {
				resp.Unauthorized("invalid api key").WithModule("AuthMW").Send(w)
				return
			}

			workspace, err := mw.workspaces.GetByID(ctx, matched.WorkspaceID)
			if err != nil {
				resp.Unauthorized("workspace not found").WithModule("AuthMW").Send(w)
				return
			}

			ctx = authz.WithWorkspace(ctx, workspace)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (mw *AuthMiddleware) AnyAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx, span := mw.tracer.Start(ctx, "Middleware.AnyAuth")
			defer span.End()

			// try API key first
			rawKey := r.Header.Get("X-API-Key")
			if rawKey != "" {
				if len(rawKey) >= 11 {
					prefix := rawKey[:11]
					candidates, err := mw.apiKeys.GetByPrefix(ctx, prefix)
					if err == nil && len(candidates) > 0 {
						for _, candidate := range candidates {
							if err := bcrypt.CompareHashAndPassword([]byte(candidate.KeyHash), []byte(rawKey)); err == nil {
								workspace, err := mw.workspaces.GetByID(ctx, candidate.WorkspaceID)
								if err == nil {
									ctx = authz.WithWorkspace(ctx, workspace)
									next.ServeHTTP(w, r.WithContext(ctx))
									return
								}
							}
						}
					}
				}
				resp.Unauthorized("invalid api key").WithModule("AuthMW").Send(w)
				return
			}

			// fall back to cookie auth
			accessTokenCookie, err := r.Cookie("access_token")
			if err != nil {
				resp.Unauthorized("missing credentials").WithModule("AuthMW").Send(w)
				return
			}

			token, err := mw.gaClient.Tokens.ValidateToken(ctx, accessTokenCookie.Value)
			if err != nil {
				resp.Unauthorized(err.Error()).WithModule("AuthMW").Send(w)
				return
			}

			subject, err := authz.GetSubjectFromToken(token)
			if err != nil {
				resp.Unauthorized(err.Error()).WithModule("AuthMW").Send(w)
				return
			}

			ctx = authz.WithSubject(ctx, subject)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
