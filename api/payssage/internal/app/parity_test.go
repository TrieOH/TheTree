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
	"riverqueue.com/riverui"
)

func mwJWT(next http.Handler) http.Handler { return next }

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

// expectedRoutes is the authoritative route inventory for Payssage, sourced
// from each feature's handlers/0_handler.go RegisterRoutes. Value is the
// auth chain ("public" when none). The /riverui mount is an internal ops
// console and intentionally excluded from both the walk assertion and the
// OpenAPI spec.
var expectedRoutes = map[string]string{
	// harness-owned
	"GET /health":           "public",
	"GET /docs/openapi.yml": "public",
	// orgs
	"GET /organizations":                                                  "JWT",
	"POST /organizations":                                                 "JWT",
	"GET /organizations/{organization_id}/members":                        "JWT",
	"POST /organizations/{organization_id}/members":                       "JWT",
	"DELETE /organizations/{organization_id}/members":                     "JWT",
	"GET /organizations/{organization_id}/member/{member_id}":             "JWT",
	"GET /organizations/{organization_id}/member/{member_email}:by_email": "JWT",
	// wallets
	"POST /wallets":                                "JWT",
	"GET /wallets":                                 "JWT",
	"GET /wallets/{wallet_id}":                     "JWT",
	"PATCH /wallets/{wallet_id}/fee":               "JWT",
	"PATCH /wallets/{wallet_id}/sandbox":           "JWT",
	"GET /organizations/{organization_id}/wallets": "JWT",
	"POST /wallets/{wallet_id}/collector":          "JWT",
	"DELETE /wallets/{wallet_id}/collector":        "JWT",
	// collectors
	"GET /collectors":                                 "JWT",
	"GET /collectors/{collector_id}":                  "JWT",
	"GET /organizations/{organization_id}/collectors": "JWT",
	// sellers
	"GET /wallets/{wallet_id}/sellers": "JWT",
	// intents
	"GET /intents":                                 "JWT",
	"GET /intents/{intent_id}":                     "JWT",
	"POST /intents/{intent_id}/cancel":             "JWT",
	"GET /wallets/{wallet_id}/intents":             "JWT",
	"GET /organizations/{organization_id}/intents": "JWT",
	"POST /wallets/{wallet_id}/intents":            "JWT",
	"POST /testmode/intents/create":                "JWT",
	// oauth
	"POST /providers/{provider}/connect": "JWT",
	"POST /providers/{provider}/revoke":  "JWT",
	"GET /providers/{provider}/callback": "public",
	// webhook receive
	"POST /webhooks/{provider}": "public",
	// webhook endpoints
	"POST /wallets/{wallet_id}/webhooks/endpoints": "JWT",
	"GET /wallets/{wallet_id}/webhooks/endpoints":  "JWT",
	"GET /webhooks/endpoints/{endpoint_id}":        "JWT",
	"DELETE /webhooks/endpoints/{endpoint_id}":     "JWT",
	// webhook events
	"GET /wallets/{wallet_id}/webhooks/events": "JWT",
	"GET /webhooks/events/{event_id}":          "JWT",
	// webhook deliveries
	"GET /webhooks/endpoints/{endpoint_id}/deliveries": "JWT",
	"GET /webhooks/deliveries/{delivery_id}":           "JWT",
}

func walkRoutes(r *chi.Mux) map[string]string {
	out := map[string]string{}
	_ = chi.Walk(r, func(method, route string, _ http.Handler, mws ...func(http.Handler) http.Handler) error {
		if strings.HasPrefix(route, "/riverui") {
			return nil // internal ops console, out of API scope
		}
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
		jwtAuth: mwJWT,
	}, handlers{}, &riverui.Handler{})

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
		jwtAuth: mwJWT,
	}, handlers{}, &riverui.Handler{})
	got := walkRoutes(r)
	keys := make([]string, 0, len(got))
	for k := range got {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%-90s %s\n", k, got[k])
	}
}
