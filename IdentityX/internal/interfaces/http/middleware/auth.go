package middleware

import (
	"IdentityX/internal/features/security"
	"IdentityX/internal/shared/authz"
	"IdentityX/internal/shared/ports"
	"context"
	"net/http"
	"strings"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/fail/v3"
	"go.opentelemetry.io/otel/trace"
)

type AuthMiddleware struct {
	authenticator security.CommandService
	tracer        trace.Tracer
	cache         ports.RedisCacheService
	issuer        string
}

func NewAuthMiddleware(
	authenticator security.CommandService,
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

// Auth is a middleware function that checks for valid access and refresh security.
// It validates the security, checks if the refresh token is revoked, and creates a principal from the security.
// The principal is then added to the request context.
func (mw *AuthMiddleware) Auth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx, span := mw.tracer.Start(ctx, "Middleware.Auth")
			defer span.End()

			in := security.AuthenticateRequestInput{
				Issuer: mw.issuer,
			}

			// ⭐ API KEY (highest priority)
			apiKey := strings.TrimSpace(r.Header.Get("X-API-Key"))
			if apiKey != "" {
				in.ApiKey = apiKey
				mw.handleAuth(ctx, w, r, next, in)
				return
			}

			// ⭐ Authorization header (Primary auth)
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				resp.Unauthorized().WithMsg("missing access token").WithModule("AuthMW").Send(w)
				return
			}

			_, tokenStr, found := strings.Cut(authHeader, "Bearer ")
			if !found || tokenStr == "" {
				resp.Unauthorized().WithMsg("invalid authorization header").WithModule("AuthMW").Send(w)
				return
			}

			in.AccessToken = tokenStr
			mw.handleAuth(ctx, w, r, next, in)
		})
	}
}

func (mw *AuthMiddleware) handleAuth(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	next http.Handler,
	in security.AuthenticateRequestInput,
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

	ctx = authz.WithPrincipal(ctx, principal)
	next.ServeHTTP(w, r.WithContext(ctx))
}
