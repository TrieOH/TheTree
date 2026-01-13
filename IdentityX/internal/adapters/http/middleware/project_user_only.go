package middleware

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/application/auth"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

func ProjectUserOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		principal, err := auth.RequirePrincipal(ctx)
		if err != nil {
			resp.FromError(apierr.FromService(nil, err)).WithModule("ProjectUserOnlyMW").Send(w)
			return
		}

		if principal.ProjectID == nil {
			err = apierr.ErrForbidden.WithMsg("only project users can access this endpoint").WithID(apierr.AuthNotProjectUser)
			resp.FromError(err).WithModule("ProjectUserOnlyMW").Send(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}
