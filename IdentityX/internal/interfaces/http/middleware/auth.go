package middleware

import (
	"IdentityX/internal/features/auth"
	"IdentityX/internal/platform/telemetry"
	authz2 "IdentityX/internal/shared/authz"
	"IdentityX/internal/shared/ports"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/fail/v3"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type AuthMiddleware struct {
	authenticator auth.CommandService
	tracer        trace.Tracer
	cache         ports.RedisCacheService
	issuer        string
}

func NewAuthMiddleware(
	authenticator auth.CommandService,
	tracer trace.Tracer,
	cache ports.RedisCacheService,
	issuer string,
) *AuthMiddleware {
	return &AuthMiddleware{
		authenticator: authenticator,
		tracer:        tracer,
		cache:         cache,
		issuer:        issuer,
	}
}

// FIXME:
// Auto-exchange flow is not yet implemented for backward compatibility.
// Current behavior requires clients to explicitly call /exchange
// when no valid svc_session cookie is present.
//
// Future design:
// Middleware should transparently perform exchange when:
//
// 1. svc_session cookie missing OR cache miss
// 2. Authorization Bearer global access token present
//
// This will enable:
// - zero-RTT relying-party session bootstrap
// - simpler frontend logic
// - better SSR / websocket auth ergonomics
// - resilience to session eviction

// Auth is a middleware function that checks for valid access and refresh tokens.
// It validates the tokens, checks if the refresh token is revoked, and creates a principal from the tokens.
// The principal is then added to the request context.
func (mw *AuthMiddleware) Auth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx, span := mw.tracer.Start(ctx, "Middleware.Auth")
			defer span.End()

			in := auth.AuthenticateRequestInput{
				Issuer: mw.issuer,
			}

			// ⭐ API KEY (highest priority)
			apiKey := strings.TrimSpace(r.Header.Get("X-API-Key"))
			if apiKey != "" {
				in.ApiKey = apiKey
				mw.handleAuth(ctx, w, r, next, in)
				return
			}

			svcCookie, err := r.Cookie("svc_session")
			if err == nil {
				if svcCookie.Value == "" {
					resp.Unauthorized().WithMsg("Empty service cookie").WithModule("AuthMW").Send(w)
					return
				}

				key := "svc_session:" + svcCookie.Value

				var snapshotAny any
				var found bool
				snapshotAny, found, err = mw.cache.GetAny(ctx, "svc_session:"+svcCookie.Value)
				if err != nil {
					// IMPORTANT:
					// Cache failure must NOT hard fail auth.
					// We fall back to Bearer / other auth sources.
				} else if found {
					data, ok := snapshotAny.([]byte)
					if !ok {
						telemetry.Log().Error("unexpected cache type", zap.Any("snapshotAny", snapshotAny))
						resp.Unauthorized().WithMsg("invalid session").Send(w)
						return
					}

					var snapshot authz2.ServiceSnapshot
					if err := json.Unmarshal(data, &snapshot); err != nil {
						telemetry.Log().Error("failed to unmarshal session", zap.Error(err))
						_ = mw.cache.Delete(ctx, key)
						resp.Unauthorized().WithMsg("invalid session").Send(w)
						return
					}

					// TTL safety guard (important)
					if time.Now().After(snapshot.AccessData.ExpiresAt.Time) {
						err = mw.cache.Delete(ctx, "svc_session:"+svcCookie.Value)
						if err != nil {
							telemetry.Log().Error("Error deleting service session", zap.Error(err))
						}
						resp.Unauthorized().WithMsg("session expired").WithModule("AuthMW").Send(w)
						return
					}

					principal := snapshot.ToPrincipal()
					ctx = authz2.WithPrincipal(ctx, principal)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
				telemetry.Log().Info("Auth cache entry not found")
			} else {
				telemetry.Log().Info("Error getting cookie", zap.Error(err))
			}

			// ⭐ Authorization header (PRIMARY auth)
			authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
			if authHeader != "" {

				if !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
					resp.BadRequest().WithMsg("invalid authorization header").WithModule("AuthMW").Send(w)
					return
				}

				token := strings.TrimSpace(authHeader[7:])
				if token == "" {
					resp.BadRequest().WithMsg("empty bearer token").WithModule("AuthMW").Send(w)
					return
				}

				in.AccessToken = token
				mw.handleAuth(ctx, w, r, next, in)
				return
			}

			// ⭐ Cookie fallback (legacy / browser only)
			cookie, err := r.Cookie("access_token")
			if err == nil {
				in.AccessToken = cookie.Value
				mw.handleAuth(ctx, w, r, next, in)
				return
			}

			resp.Unauthorized().WithMsg("missing authentication").WithModule("AuthMW").Send(w)
		})
	}
}

func (mw *AuthMiddleware) handleAuth(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	next http.Handler,
	in auth.AuthenticateRequestInput,
) {
	principal, err := mw.authenticator.AuthenticateRequest(ctx, in)
	if err != nil {
		rs, convErr := fail.ToAs[*resp.Response](fail.AsFail(err), "http")
		if convErr != nil {
			resp.InternalServerError().WithModule("AuthMW").Send(w)
			return
		}
		rs.WithModule("AuthMW").Send(w)
		return
	}

	ctx = authz2.WithPrincipal(ctx, principal)
	next.ServeHTTP(w, r.WithContext(ctx))
}
