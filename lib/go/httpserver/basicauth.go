package httpserver

import (
	"crypto/subtle"
	"net/http"
	"os"
)

const (
	basicAuthUserEnv = "SIMPLE_AUTH_USER"
	basicAuthPassEnv = "SIMPLE_AUTH_PASS"
)

// BasicAuth protects an HTTP subtree with HTTP Basic auth backed by the
// SIMPLE_AUTH_USER / SIMPLE_AUTH_PASS env vars. It is fail-closed: when
// either variable is unset the subtree answers 503 instead of becoming
// public, so a forgotten deployment cannot silently widen the attack
// surface. Credentials are compared in constant time.
//
// Backends that mount an operational dashboard (e.g. riverui) apply this
// middleware inside the mount's chi group: the spec-derived auth chains
// cover only generated operations, so non-spec routes would otherwise be
// served without authentication.
func BasicAuth(next http.Handler) http.Handler {
	user := os.Getenv(basicAuthUserEnv)
	pass := os.Getenv(basicAuthPassEnv)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user == "" || pass == "" {
			http.Error(w, "basic auth not configured", http.StatusServiceUnavailable)
			return
		}

		reqUser, reqPass, ok := r.BasicAuth()
		if !ok ||
			subtle.ConstantTimeCompare([]byte(reqUser), []byte(user)) != 1 ||
			subtle.ConstantTimeCompare([]byte(reqPass), []byte(pass)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
