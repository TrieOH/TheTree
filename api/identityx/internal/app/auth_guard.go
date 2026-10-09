package app

import (
	"net/http"

	"IdentityX/internal/setup"

	"github.com/MintzyG/fun"
)

// setupGuard returns the middleware that gates every operation (except the
// two /auth/setup routes, which check the state themselves) until setup has
// completed.
func setupGuard() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			done, err := setup.Check(r.Context())
			if err != nil {
				fun.ServiceUnavailable("could not check setup state").Send(w)
				return
			}
			if !done {
				fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
