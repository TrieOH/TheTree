package middleware

import (
	"GoAuth/internal/adapters/http/utils"
	"GoAuth/internal/apierr"
	"GoAuth/internal/application/auth"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/fail"
)

// ClientOnly is a middleware that ensures that the request is made by a client and not a project user.
func ClientOnly() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			principal, err := auth.RequirePrincipal(ctx)
			if err != nil {
				resp.FromError(apierr.FromService(nil, err)).WithModule("ClientOnlyMW").Send(w)
				return
			}

			if principal.ProjectID != nil {
				rs, ok := utils.Sender(fail.New(apierr.AuthNotClient), "ClientOnlyMW", w)
				if ok {
					rs.Send(w)
				}
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
