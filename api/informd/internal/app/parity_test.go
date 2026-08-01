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

func mwJWT(next http.Handler) http.Handler     { return next }
func mwAnyAuth(next http.Handler) http.Handler { return next }

func mwName(mw func(http.Handler) http.Handler) string {
	fn := runtime.FuncForPC(reflect.ValueOf(mw).Pointer())
	name := fn.Name()
	if i := strings.LastIndex(name, "."); i >= 0 {
		name = name[i+1:]
	}
	if strings.HasPrefix(name, "func") {
		return ""
	}
	return strings.TrimPrefix(name, "mw")
}

// expectedRoutes is the authoritative route inventory for Informd, sourced
// from each feature's handlers/0_handler.go RegisterRoutes. Value is the
// auth chain ("public" when none).
var expectedRoutes = map[string]string{
	// harness-owned
	"GET /health":           "public",
	"GET /docs/openapi.yml": "public",
	// forms (note: /forms/{form_id}/asnwerable is intentionally public and
	// keeps the original spelling used by the frontend)
	"GET /forms/{form_id}/asnwerable":      "public",
	"POST /forms":                          "AnyAuth",
	"GET /forms":                           "AnyAuth",
	"GET /forms/archived":                  "AnyAuth",
	"GET /forms/{form_id}/full":            "AnyAuth",
	"GET /forms/{form_id}/members":         "AnyAuth",
	"POST /forms/{form_id}/members":        "AnyAuth",
	"DELETE /forms/{form_id}/members":      "AnyAuth",
	"POST /forms/{form_id}/open":           "AnyAuth",
	"POST /forms/{form_id}/close":          "AnyAuth",
	"POST /forms/{form_id}/archive":        "AnyAuth",
	"POST /forms/{form_id}/redraft":        "AnyAuth",
	"GET /forms/{form_id}/responses/count": "AnyAuth",
	// namespaces
	"GET /namespaces":                                                "JWT",
	"POST /namespaces":                                               "JWT",
	"GET /namespaces/{namespace_id}/members":                         "JWT",
	"POST /namespaces/{namespace_id}/members":                        "JWT",
	"DELETE /namespaces/{namespace_id}/members":                      "JWT",
	"POST /namespaces/{namespace_id}/forms":                          "JWT",
	"GET /namespaces/{namespace_id}/forms":                           "JWT",
	"GET /namespaces/{namespace_id}/forms/archived":                  "JWT",
	"GET /namespaces/{namespace_id}/forms/{form_id}/full":            "JWT",
	"GET /namespaces/{namespace_id}/forms/{form_id}/members":         "JWT",
	"POST /namespaces/{namespace_id}/forms/{form_id}/members":        "JWT",
	"DELETE /namespaces/{namespace_id}/forms/{form_id}/members":      "JWT",
	"POST /namespaces/{namespace_id}/forms/{form_id}/open":           "JWT",
	"POST /namespaces/{namespace_id}/forms/{form_id}/close":          "JWT",
	"POST /namespaces/{namespace_id}/forms/{form_id}/archive":        "JWT",
	"POST /namespaces/{namespace_id}/forms/{form_id}/redraft":        "JWT",
	"GET /namespaces/{namespace_id}/forms/{form_id}/responses/count": "JWT",
	// steps
	"POST /forms/{form_id}/steps":                           "AnyAuth",
	"PUT /forms/{form_id}/steps":                            "AnyAuth",
	"GET /forms/{form_id}/steps":                            "AnyAuth",
	"GET /namespaces/{namespace_id}/forms/{form_id}/steps":  "AnyAuth",
	"POST /namespaces/{namespace_id}/forms/{form_id}/steps": "AnyAuth",
	"PUT /namespaces/{namespace_id}/forms/{form_id}/steps":  "AnyAuth",
	// fields
	"POST /forms/{form_id}/steps/{step_id}/fields":                                            "AnyAuth",
	"PUT /forms/{form_id}/steps/{step_id}/fields":                                             "AnyAuth",
	"GET /forms/{form_id}/steps/{step_id}/fields":                                             "AnyAuth",
	"GET /forms/{form_id}/steps/{step_id}/fields/{field_id}/select":                           "AnyAuth",
	"DELETE /forms/{form_id}/steps/{step_id}/fields/{field_id}":                               "AnyAuth",
	"PUT /forms/{form_id}/steps/{step_id}/fields/{field_id}/select":                           "AnyAuth",
	"POST /namespaces/{namespace_id}/forms/{form_id}/steps/{step_id}/fields":                  "AnyAuth",
	"PUT /namespaces/{namespace_id}/forms/{form_id}/steps/{step_id}/fields":                   "AnyAuth",
	"GET /namespaces/{namespace_id}/forms/{form_id}/steps/{step_id}/fields":                   "AnyAuth",
	"GET /namespaces/{namespace_id}/forms/{form_id}/steps/{step_id}/fields/{field_id}/select": "AnyAuth",
	"DELETE /namespaces/{namespace_id}/forms/{form_id}/steps/{step_id}/fields/{field_id}":     "AnyAuth",
	"PUT /namespaces/{namespace_id}/forms/{form_id}/steps/{step_id}/fields/{field_id}/select": "AnyAuth",
	// responses
	"POST /forms/{form_id}/responses": "public",
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

func TestRouterParity(t *testing.T) {
	r := chi.NewRouter()
	// /health is harness-owned (httpserver.NewRouter); mirror its registration.
	r.Get("/health", http.NotFoundHandler().ServeHTTP)
	// /docs/openapi.yml is harness-owned (httpserver.NewRouter); mirror its registration.
	r.Get("/docs/openapi.yml", http.NotFoundHandler().ServeHTTP)
	registerRoutes(r, middlewares{
		jwt:     mwJWT,
		anyAuth: mwAnyAuth,
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
		jwt:     mwJWT,
		anyAuth: mwAnyAuth,
	}, handlers{})
	got := walkRoutes(r)
	keys := make([]string, 0, len(got))
	for k := range got {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%-95s %s\n", k, got[k])
	}
}
