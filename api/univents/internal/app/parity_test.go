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

// expectedRoutes is the authoritative route inventory for Univents, sourced
// from each feature's handlers/0_handler.go RegisterRoutes. Value is the
// auth chain ("public" when none). The /riverui mount is an internal ops
// console and intentionally excluded from both the walk assertion and the
// OpenAPI spec.
var expectedRoutes = map[string]string{
	// harness-owned
	"GET /health":           "public",
	"GET /docs/openapi.yml": "public",
	// events
	"GET /events":                                 "public",
	"GET /events/{event_slug}:by-slug":            "public",
	"POST /events":                                "JWT",
	"GET /events/owned":                           "JWT",
	"GET /events/joined":                          "JWT",
	"PATCH /events/{event_id}":                    "JWT",
	"POST /events/{event_id}/publish":             "JWT",
	"POST /events/{event_id}/discontinue":         "JWT",
	"GET /events/{event_id}/members":              "JWT",
	"POST /events/{event_id}/members":             "JWT",
	"DELETE /events/{event_id}/members/{user_id}": "JWT",
	// editions
	"GET /events/{event_slug}:by-slug/editions/{edition_slug}:by-slug": "public",
	"GET /events/{event_id}/editions":                                  "public",
	"GET /events/{event_id}/editions/active":                           "public",
	"GET /events/{event_id}/editions/past":                             "public",
	"GET /events/{event_id}/editions/upcoming":                         "public",
	"POST /events/{event_id}/editions":                                 "JWT",
	"GET /events/{event_id}/editions/draft":                            "JWT",
	"PATCH /events/{event_id}/editions/{edition_id}":                   "JWT",
	"POST /events/{event_id}/editions/{edition_id}/publish":            "JWT",
	// ticket_types
	"GET /editions/{edition_id}/ticket-types":  "public",
	"GET /ticket-types/{ticket_type_id}":       "public",
	"PATCH /ticket-types/{ticket_type_id}":     "JWT",
	"POST /editions/{edition_id}/ticket-types": "JWT",
	// products
	"GET /editions/{edition_id}/products":                       "public",
	"GET /editions/{edition_id}/products/{vendor_code}:by-code": "public",
	"GET /editions/{edition_id}/variants/{vendor_code}:by-code": "public",
	"GET /products/{product_id}":                                "public",
	"GET /products/{product_id}/variants":                       "public",
	"POST /editions/{edition_id}/products":                      "JWT",
	"POST /products/{product_id}/variants":                      "JWT",
	"PATCH /products/{product_id}":                              "JWT",
	"PATCH /variants/{variant_id}":                              "JWT",
	"DELETE /products/{product_id}":                             "JWT",
	"DELETE /variants/{variant_id}":                             "JWT",
	// programs
	"GET /editions/{edition_id}/programs":     "public",
	"POST /editions/{edition_id}/programs":    "JWT",
	"GET /programs/{program_id}":              "public",
	"PATCH /programs/{program_id}":            "JWT",
	"DELETE /programs/{program_id}":           "JWT",
	"GET /programs/{program_id}/occurrences":  "public",
	"POST /programs/{program_id}/occurrences": "JWT",
	"GET /editions/{edition_id}/occurrences":  "public",
	"GET /occurrences/{occurrence_id}":        "public",
	"PATCH /occurrences/{occurrence_id}":      "JWT",
	"DELETE /occurrences/{occurrence_id}":     "JWT",
	// badges
	"POST /editions/{edition_id}/badges": "JWT",
	"GET /editions/{edition_id}/badges":  "JWT",
	"GET /badges/{template_id}":          "JWT",
	"DELETE /badges/{template_id}":       "JWT",
	// signatures
	"GET /editions/{edition_id}/signatures":          "public",
	"GET /signatures/{signature_id}":                 "public",
	"POST /editions/{edition_id}/signatures":         "JWT",
	"DELETE /signatures/{signature_id}":              "JWT",
	"GET /editions/{edition_id}/signature-requests":  "public",
	"GET /signature-requests/{request_id}":           "public",
	"POST /editions/{edition_id}/signature-requests": "JWT",
	"POST /signature-requests/fulfill":               "public",
	"POST /signature-requests/deny":                  "public",
	"POST /signature-requests/{request_id}/cancel":   "JWT",
	"POST /signatures/revoke":                        "public",
	// certifications
	"GET /editions/{edition_id}/certifications/templates":       "public",
	"GET /certifications/templates/{template_id}":               "public",
	"POST /editions/{edition_id}/certifications/templates":      "JWT",
	"PUT /certifications/templates/{template_id}":               "JWT",
	"DELETE /certifications/templates/{template_id}":            "JWT",
	"POST /certifications/templates/{template_id}/link":         "JWT",
	"DELETE /certifications/templates/{template_id}/link":       "JWT",
	"GET /certifications/templates/{template_id}/links":         "public",
	"GET /verify/{hash}":                                        "public",
	"GET /certifications/{cert_id}":                             "JWT",
	"GET /editions/{edition_id}/certifications":                 "JWT",
	"GET /certifications":                                       "JWT",
	"POST /certifications/{cert_id}/invalidate":                 "JWT",
	"GET /editions/{edition_id}/certifications/emission-errors": "JWT",
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
		jwt: mwJWT,
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
		jwt: mwJWT,
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
