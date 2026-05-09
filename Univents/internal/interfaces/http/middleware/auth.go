package middleware

import (
	"net/http"
	"strings"
	"univents/internal/shared/authz"

	"git.trieoh.com/TrieOH/IdentityX-SDK-Go"
	"github.com/MintzyG/fun"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type AuthMiddleware struct {
	idxClient idx.Client
	tracer    trace.Tracer
}

func NewAuthMiddleware(
	idxClient *idx.Client,
	tracer trace.Tracer,
) *AuthMiddleware {
	return &AuthMiddleware{
		idxClient: *idxClient,
		tracer:    tracer,
	}
}

// Auth is a middleware that validates the Authorization header Bearer token.
// It injects the subject into the request context if valid.
func (mw *AuthMiddleware) Auth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx, span := mw.tracer.Start(ctx, "Middleware.Auth")
			defer span.End()

			var err error
			defer func() {
				span.SetAttributes(attribute.Bool("success", err == nil))
			}()

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				fun.Unauthorized().WithMsg("missing access token").WithModule("AuthMW").Send(w)
				return
			}

			_, tokenStr, found := strings.Cut(authHeader, "Bearer ")
			if !found || tokenStr == "" {
				fun.Unauthorized().WithMsg("invalid authorization header").WithModule("AuthMW").Send(w)
				return
			}

			accessClaims, err := mw.idxClient.Tokens.VerifyAccessToken(ctx, tokenStr)
			if err != nil {
				fun.Unauthorized().WithMsg("invalid access token").WithModule("AuthMW").Send(w)
				return
			}

			// Inject subject into context
			subject := authz.UserSubject{
				ID:    accessClaims.Sub.ID,
				Email: accessClaims.Sub.Email,
			}
			ctx = authz.WithSubject(ctx, &subject)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
