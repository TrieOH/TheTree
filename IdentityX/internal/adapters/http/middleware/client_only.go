package middleware

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/application/auth"
	"net/http"
)

// ClientOnly is a middleware that ensures that the request is made by a client and not a project user.
func ClientOnly() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			principal, err := auth.RequirePrincipal(ctx)
			if err != nil {
				ErrToResp(err).WithModule("ClientOnlyMW").Send(w)
				return
			}

			if principal.ProjectID != nil {
				mwErr := apierr.ErrForbidden.WithMsg("only clients can access this endpoint").WithID(apierr.AuthNotClient)
				ErrToResp(mwErr).WithModule("ClientOnlyMW").Send(w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
