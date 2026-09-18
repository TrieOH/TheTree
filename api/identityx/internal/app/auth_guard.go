package app

import (
	"net/http"

	"lib/globals"

	"github.com/MintzyG/fun"
)

// setupGuard returns the middleware that gates every operation (except the
// two /auth/setup routes, which manage the flag themselves) until setup has
// completed.
func setupGuard() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !globals.SetupComplete() {
				fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
