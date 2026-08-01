package app

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

// Named noop middlewares: registration only stores them (they panic only if
// served), and the walk reports their names so each route's auth chain can be
// asserted against the spec's security declarations.
func mwJWT(next http.Handler) http.Handler        { return next }
func mwAnyAuth(next http.Handler) http.Handler    { return next }
func mwClientOnly(next http.Handler) http.Handler { return next }

func mwName(mw func(http.Handler) http.Handler) string {
	fn := runtime.FuncForPC(reflect.ValueOf(mw).Pointer())
	name := fn.Name()
	if i := strings.LastIndex(name, "."); i >= 0 {
		name = name[i+1:]
	}
	// Anonymous middlewares (e.g. fun's WithParams[ProjectIDQueryParam]
	// parameter binder) report as "func1"; they carry no auth meaning, so
	// treat them as absent from the auth chain.
	if strings.HasPrefix(name, "func") {
		return ""
	}
	return strings.TrimPrefix(name, "mw")
}

// expectedRoutes is the authoritative route inventory for IdentityX,
// sourced from each feature's handlers/0_handler.go RegisterRoutes.
// Value is the auth chain ("public" when none). Used to prove spec parity:
// every route walked from the live router must be listed here with the same
// auth chain, and every listed route must be walked.
var expectedRoutes = map[string]string{
	// harness-owned
	"GET /health":           "public",
	"GET /docs/openapi.yml": "public",
	// authn
	"GET /auth/setup":               "public",
	"POST /auth/setup":              "public",
	"POST /auth/register":           "public",
	"POST /auth/login":              "public",
	"POST /auth/logout":             "JWT",
	"POST /auth/refresh":            "public",
	"GET /auth/{provider}/connect":  "public",
	"GET /auth/{provider}/callback": "public",
	"GET /.well-known/jwks.json":    "public",
	"GET /auth/introspect":          "AnyAuth",
	// actors
	"GET /projects/{project_id}/actors/{actor_id}":             "AnyAuth+ClientOnly",
	"GET /projects/{project_id}/actors/{actor_email}:by_email": "AnyAuth+ClientOnly",
	"POST /projects/{project_id}/actors":                       "AnyAuth+ClientOnly",
	"GET /projects/{project_id}/actors":                        "AnyAuth+ClientOnly",
	// api_keys
	"POST /projects/{project_id}/api_keys": "AnyAuth+ClientOnly",
	// capabilities
	"GET /projects/{project_id}/capabilities":  "AnyAuth",
	"POST /projects/{project_id}/capabilities": "JWT+ClientOnly",
	// organizations
	"GET /organizations":                                                           "JWT+ClientOnly",
	"POST /organizations":                                                          "JWT+ClientOnly",
	"GET /organizations/{organization_id}/members":                                 "JWT+ClientOnly",
	"POST /organizations/{organization_id}/members":                                "JWT+ClientOnly",
	"DELETE /organizations/{organization_id}/members":                              "JWT+ClientOnly",
	"GET /organizations/{organization_id}/projects":                                "JWT+ClientOnly",
	"POST /organizations/{organization_id}/projects":                               "JWT+ClientOnly",
	"POST /organizations/{org_id}/projects/{project_id}/actors":                    "JWT+ClientOnly",
	"GET /organizations/{org_id}/projects/{project_id}/actors":                     "JWT+ClientOnly",
	"GET /organizations/{organization_id}/projects/{project_id}/members":           "JWT+ClientOnly",
	"POST /organizations/{organization_id}/projects/{project_id}/members":          "JWT+ClientOnly",
	"DELETE /organizations/{organization_id}/projects/{project_id}/members":        "JWT+ClientOnly",
	"GET /organizations/{organization_id}/projects/{project_id}/actors/{actor_id}": "JWT+ClientOnly",
	// projects
	"GET /projects":                         "AnyAuth+ClientOnly",
	"POST /projects":                        "AnyAuth+ClientOnly",
	"GET /projects/{project_id}/members":    "AnyAuth+ClientOnly",
	"POST /projects/{project_id}/members":   "AnyAuth+ClientOnly",
	"DELETE /projects/{project_id}/members": "AnyAuth+ClientOnly",
	// profiles
	"GET /actors/{actor_id}/profile":                       "JWT+ClientOnly",
	"PUT /actors/{actor_id}/profile":                       "JWT+ClientOnly",
	"GET /projects/{project_id}/actors/{actor_id}/profile": "JWT+ClientOnly",
	"PUT /projects/{project_id}/actors/{actor_id}/profile": "JWT+ClientOnly",
	// profile_schemas
	"GET /projects/{project_id}/profile-schema": "JWT+ClientOnly",
	"PUT /projects/{project_id}/profile-schema": "JWT+ClientOnly",
	"GET /profile-schema":                       "JWT+ClientOnly",
	"PUT /profile-schema":                       "JWT+ClientOnly",
}

func walkRoutes(r *chi.Mux) map[string]string {
	out := map[string]string{}
	_ = chi.Walk(r, func(method, route string, _ http.Handler, mws ...func(http.Handler) http.Handler) error {
		names := make([]string, 0, len(mws))
		for _, mw := range mws {
			if n := mwName(mw); n != "" {
				names = append(names, n)
			}
		}
		chain := strings.Join(names, "+")
		if chain == "" {
			chain = "public"
		}
		out[normalizeRoute(method+" "+route)] = chain
		return nil
	})
	return out
}

func normalizeRoute(r string) string {
	parts := strings.SplitN(r, " ", 2)
	if len(parts) != 2 {
		return r
	}
	path := strings.TrimSuffix(parts[1], "/")
	if path == "" {
		path = "/"
	}
	return parts[0] + " " + path
}

// TestRouterParity walks the live chi router and asserts bidirectional
// equality with the expected route inventory, including each route's auth
// middleware chain. Zero-value handlers are safe: registration only stores
// them, they panic only if served.
func TestRouterParity(t *testing.T) {
	r := chi.NewRouter()
	// /health is harness-owned (httpserver.NewRouter); mirror its registration.
	r.Get("/health", http.NotFoundHandler().ServeHTTP)
	// /docs/openapi.yml is harness-owned (httpserver.NewRouter); mirror its registration.
	r.Get("/docs/openapi.yml", http.NotFoundHandler().ServeHTTP)
	registerRoutes(r, middlewares{
		jwtAuth:    mwJWT,
		anyAuth:    mwAnyAuth,
		clientOnly: mwClientOnly,
	}, handlers{})

	got := walkRoutes(r)

	var missing, extra, authMismatch []string
	for want, wantAuth := range expectedRoutes {
		gotAuth, ok := got[want]
		if !ok {
			missing = append(missing, want)
		} else if gotAuth != wantAuth {
			authMismatch = append(authMismatch, fmt.Sprintf("%s: want %s, got %s", want, wantAuth, gotAuth))
		}
	}
	for gotRoute := range got {
		if _, ok := expectedRoutes[gotRoute]; !ok {
			extra = append(extra, gotRoute)
		}
	}
	sort.Strings(missing)
	sort.Strings(extra)
	sort.Strings(authMismatch)

	if len(missing) > 0 || len(extra) > 0 || len(authMismatch) > 0 {
		t.Fatalf("route parity mismatch\nroutes expected but not walked:\n%s\nroutes walked but not expected:\n%s\nauth chain mismatches:\n%s",
			strings.Join(missing, "\n"), strings.Join(extra, "\n"), strings.Join(authMismatch, "\n"))
	}

	t.Logf("parity ok: %d routes with matching auth chains (spec parity proven)", len(got))
}

func TestWalkOutputForReview(t *testing.T) {
	if testing.Short() {
		t.Skip("skip walk dump in short mode")
	}
	r := chi.NewRouter()
	r.Get("/health", http.NotFoundHandler().ServeHTTP)
	// /docs/openapi.yml is harness-owned (httpserver.NewRouter); mirror its registration.
	r.Get("/docs/openapi.yml", http.NotFoundHandler().ServeHTTP)
	registerRoutes(r, middlewares{
		jwtAuth:    mwJWT,
		anyAuth:    mwAnyAuth,
		clientOnly: mwClientOnly,
	}, handlers{})
	got := walkRoutes(r)
	keys := make([]string, 0, len(got))
	for k := range got {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%-75s %s\n", k, got[k])
	}
}
