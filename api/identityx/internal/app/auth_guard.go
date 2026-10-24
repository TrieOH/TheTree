package app

import (
	"net/http"

	"IdentityX/internal/setup"
	libauthz "lib/authz"

	"github.com/MintzyG/fun"
)

// chainOptions puts the setup guard in front of every operation except the
// two /auth/setup routes, which check the state themselves. The resolver
// validates the skipped operation IDs against the spec, so a renamed route
// fails boot.
func chainOptions(st *setup.State) libauthz.Options {
	return libauthz.Options{
		SetupGuard:     setupGuard(st),
		SkipSetupGuard: []string{"getSetup", "postSetup"},
	}
}

// setupGuard refuses every request until setup has completed.
func setupGuard(st *setup.State) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			done, err := st.Check(r.Context())
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
