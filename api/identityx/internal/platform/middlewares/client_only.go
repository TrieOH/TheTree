package middlewares

import (
	"net/http"

	"IdentityX/internal/shared/authz"

	"github.com/MintzyG/fun"
)

// ClientOnly is a middleware that ensures that the request is made by a client and not a project user.
func ClientOnly() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			principal, err := authz.RequirePrincipal(r.Context())
			if fun.Bail(w, err) {
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
