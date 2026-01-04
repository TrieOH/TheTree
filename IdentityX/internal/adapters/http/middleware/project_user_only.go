package middleware

import (
	"GoAuth/internal/application/authz"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

func ProjectUserOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		principal, err := authz.RequirePrincipal(ctx)
		if err != nil {
			ErrToResp(err).WithModule("ProjectUserOnlyMW").Send(w)
			return
		}

		if principal.ProjectID == nil {
			resp.Unauthorized("only project users can access this endpoint").WithModule("ProjectUserOnlyMW").Send(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}
