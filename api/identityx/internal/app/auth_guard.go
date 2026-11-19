package app

import (
	"fmt"
	"net/http"

	"IdentityX/errs"
	"IdentityX/internal/setup"
	"lib/httpserver"
	"lib/openapispec"
)

// setupOperations are the two /auth/setup routes, which check the setup
// state themselves.
var setupOperations = []string{"getSetup", "postSetup"}

// guardSetup puts the setup guard in front of every operation's chain
// except the setup routes. The guard runs first, before authentication. A
// renamed setup route fails boot.
func guardSetup(chains map[string][]func(http.Handler) http.Handler, st *setup.State) error {
	skip := make(map[string]bool, len(setupOperations))
	for _, id := range setupOperations {
		generated := openapispec.GeneratedID(id)
		if _, ok := chains[generated]; !ok {
			return fmt.Errorf("setup guard: operation %q not found in the spec", id)
		}
		skip[generated] = true
	}
	guard := setupGuard(st)
	for id, chain := range chains {
		if !skip[id] {
			chains[id] = append([]func(http.Handler) http.Handler{guard}, chain...)
		}
	}
	return nil
}

// setupGuard refuses every request until setup has completed.
func setupGuard(st *setup.State) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			done, err := st.Check(r.Context())
			if err != nil {
				httpserver.WriteError(w, r, fmt.Errorf("%w: %w", errs.ErrSetupCheckFailed, err))
				return
			}
			if !done {
				httpserver.WriteError(w, r, errs.ErrSetupRequired)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
