package middlewares

import (
	"net/http"

	"IdentityX/internal/shared/authz"

	"github.com/MintzyG/fun"
)

func ProjectUserOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		principal, err := authz.RequirePrincipal(r.Context())
		if err != nil {
			fun.Error(err).WithModule("ProjectUserOnlyMW").Send(w)
			return
		}

		if principal.ProjectID == nil {
			fun.Forbidden("only project users can access this endpoint").WithModule("ProjectUserOnlyMW").Send(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}
