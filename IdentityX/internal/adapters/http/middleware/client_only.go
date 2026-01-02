package middleware

import (
	"GoAuth/internal/application/authz"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

// ClientOnly is a middleware that ensures that the request is made by a client and not a project user.
func ClientOnly() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			principal, err := authz.RequirePrincipal(ctx)
			if err != nil {
				ErrToResp(err).WithModule("ClientOnlyMW").Send(w)
				return
			}

			if principal.ProjectID != nil {
				resp.Unauthorized("only clients can access this endpoint").WithModule("ClientOnlyMW").Send(w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
