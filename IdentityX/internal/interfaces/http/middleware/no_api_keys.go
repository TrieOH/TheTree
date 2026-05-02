package middleware

import (
	"IdentityX/internal/shared/authz"
	"net/http"

	"github.com/MintzyG/fun"
	"go.opentelemetry.io/otel/trace"
)

// NoApiKeys is a middleware that ensures that the request is NOT made using an API key.
func NoApiKeys() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx, span := MwTracer.Start(ctx, "NoApiKeys")
			trace.ContextWithSpan(ctx, span)
			defer span.End()

			principal, err := authz.RequirePrincipal(ctx)
			if err != nil {
				fun.Error(err).WithModule("NoApiKeysMW").Send(w)
				return
			}

			if principal.Method == authz.AuthMethodApiKey {
				fun.Error(fun.ErrBadRequest("api keys are not allowed for this endpoint")).WithModule("NoApiKeysMW").Send(w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
