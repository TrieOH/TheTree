package middleware

import (
	"IdentityX/internal/shared/authz"
	"net/http"

	"github.com/MintzyG/fun"
)

func ProjectUserOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, span := MwTracer.Start(ctx, "ProjectUserOnly")
		defer span.End()
		principal, err := authz.RequirePrincipal(ctx)
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
