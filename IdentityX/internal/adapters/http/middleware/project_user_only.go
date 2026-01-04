package middleware

import (
	"GoAuth/internal/application/authz"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

func ProjectUserOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		principal, err := authz.RequirePrincipal(ctx)
		if err != nil {
			ErrToResp(err).WithModule("ClientOnlyMW").Send(w)
			return
		}

		if principal.ProjectID == nil {
			resp.Unauthorized("only project users can access this endpoint").WithModule("ProjectUserOnlyMW").Send(w)
			return
		}

		h.ServeHTTP(w, r.WithContext(r.Context()))
	}
}
