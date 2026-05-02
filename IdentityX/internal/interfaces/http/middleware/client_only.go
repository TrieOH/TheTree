package middleware

import (
	"IdentityX/internal/shared/authz"
	"net/http"

	"github.com/MintzyG/fun"
	"go.opentelemetry.io/otel/trace"
)

// ClientOnly is a middleware that ensures that the request is made by a client and not a project user.
func ClientOnly() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx, span := MwTracer.Start(ctx, "ClientOnly")
			trace.ContextWithSpan(ctx, span)
			defer span.End()
			principal, err := authz.RequirePrincipal(ctx)
			if err != nil {
				fun.Error(err).WithModule("ClientOnlyMW").Send(w)
				return
			}

			if principal.Method == authz.AuthMethodSession && principal.ProjectID != nil {
				fun.Forbidden("only clients can access this endpoint").WithModule("ClientOnlyMW").Send(w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
