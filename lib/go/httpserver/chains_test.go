package httpserver

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"lib/httpcache"
	"lib/openapispec"

	fp "github.com/MintzyG/problematic"
)

var invocations []string

func named(name string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			invocations = append(invocations, name)
			next.ServeHTTP(w, r)
		})
	}
}

// run executes a resolved chain and returns the middleware names that ran,
// in order.
func run(chain []func(http.Handler) http.Handler) []string {
	invocations = nil
	var next http.Handler = http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	for i := range slices.Backward(chain) {
		next = chain[i](next)
	}
	next.ServeHTTP(httptest.NewRecorder(), httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil))
	return invocations
}

const (
	specWithDefault = `
openapi: 3.1.0
info: {title: t, version: "1"}
security:
  - bearerAuth: []
paths:
  /public:
    get:
      operationId: getPublic
      security: []
  /defaulted:
    get:
      operationId: getDefaulted
  /jwt:
    get:
      operationId: getJWT
      security:
        - bearerAuth: []
  /any:
    get:
      operationId: getAny
      security:
        - bearerAuth: []
        - apiKeyAuth: []
  /apikey:
    get:
      operationId: getAPIKey
      security:
        - apiKeyAuth: []
  /setup:
    get:
      operationId: getSetup
      security: []
`
	specNoDefault = `
openapi: 3.1.0
info: {title: t, version: "1"}
paths:
  /x:
    get:
      operationId: getX
`
	specUnknownScheme = `
openapi: 3.1.0
info: {title: t, version: "1"}
paths:
  /y:
    get:
      operationId: getY
      security:
        - oauth2: []
`
	specAndBlock = `
openapi: 3.1.0
info: {title: t, version: "1"}
paths:
  /z:
    get:
      operationId: getZ
      security:
        - bearerAuth: []
          apiKeyAuth: []
`
	specWithScope = `
openapi: 3.1.0
info: {title: t, version: "1"}
security:
  - bearerAuth: []
paths:
  /scoped:
    get:
      operationId: getScoped
      x-scope: platform-only
  /unscoped:
    get:
      operationId: getUnscoped
`

	specScopePublic = `
openapi: 3.1.0
info: {title: t, version: "1"}
paths:
  /scoped-public:
    get:
      operationId: getScopedPublic
      security: []
      x-scope: platform-only
`

	specScopeUnknown = `
openapi: 3.1.0
info: {title: t, version: "1"}
security:
  - bearerAuth: []
paths:
  /scoped-unknown:
    get:
      operationId: getScopedUnknown
      x-scope: nobody-registered-this
`
)

var testPrimitives = Primitives{
	JWT:    named("JWT"),
	APIKey: named("APIKey"),
	Any:    named("AnyAuth"),
}

// scopeNamed records its invocation into the shared invocations slice, so
// chain-ordering tests can see the scope middleware run in place.
func scopeNamed() ScopeChecker {
	return func(context.Context) error {
		invocations = append(invocations, "scope")
		return nil
	}
}

// errScopeDenied is a scope checker's rejection; the test mapper answers it
// with a 403.
var errScopeDenied = errors.New("scope denied")

// rejectingScope rejects every call, so the middleware path (write the
// error, stop the chain) can be pinned.
func rejectingScope(context.Context) error {
	return errScopeDenied
}

func chainsFor(spec string, primitives Primitives) (map[string][]func(http.Handler) http.Handler, error) {
	ops, err := openapispec.Operations([]byte(spec))
	if err != nil {
		return nil, err
	}
	return Chains(ops, primitives)
}

func mustChains(t *testing.T, spec string, primitives Primitives) map[string][]func(http.Handler) http.Handler {
	t.Helper()
	chains, err := chainsFor(spec, primitives)
	if err != nil {
		t.Fatalf("Chains: %v", err)
	}
	return chains
}

func TestChainsDefaultSecurityAndOverrides(t *testing.T) {
	r := mustChains(t, specWithDefault, testPrimitives)

	if got := run(r["GetPublic"]); len(got) != 0 {
		t.Fatalf("getPublic must be public, ran %v", got)
	}
	if got := run(r["GetDefaulted"]); len(got) != 1 || got[0] != "JWT" {
		t.Fatalf("getDefaulted must inherit default bearerAuth, ran %v", got)
	}
	if got := run(r["GetJWT"]); len(got) != 1 || got[0] != "JWT" {
		t.Fatalf("getJWT must be bearerAuth, ran %v", got)
	}
	if got := run(r["GetAPIKey"]); len(got) != 1 || got[0] != "APIKey" {
		t.Fatalf("getAPIKey must be apiKeyAuth, ran %v", got)
	}
	if got := run(r["GetAny"]); len(got) != 1 || got[0] != "AnyAuth" {
		t.Fatalf("getAny must derive the OR combination to AnyAuth, ran %v", got)
	}
}

func TestChainsNoDefaultMeansPublic(t *testing.T) {
	r := mustChains(t, specNoDefault, testPrimitives)
	if got := run(r["GetX"]); len(got) != 0 {
		t.Fatalf("operation without default security must be public, ran %v", got)
	}
}

func TestChainsUnregisteredCombinationFails(t *testing.T) {
	_, err := chainsFor(specUnknownScheme, testPrimitives)
	if err == nil {
		t.Fatal("want error for unregistered security combination, got nil")
	}
}

func TestChainsAndBlockFails(t *testing.T) {
	_, err := chainsFor(specAndBlock, testPrimitives)
	if err == nil || !strings.Contains(err.Error(), "separate blocks") {
		t.Fatalf("want error for single-block AND security, got %v", err)
	}
}

func TestChainsNilPrimitiveFails(t *testing.T) {
	_, err := chainsFor(specWithDefault, Primitives{})
	if err == nil || !strings.Contains(err.Error(), "is nil") {
		t.Fatalf("want error for nil primitive, got %v", err)
	}
}

func TestChainsScopeRunsAfterAuthn(t *testing.T) {
	r := mustChains(t, specWithScope, Primitives{
		JWT:    named("JWT"),
		Scopes: map[string]ScopeChecker{"platform-only": scopeNamed()},
	})

	if got := run(r["GetScoped"]); len(got) != 2 || got[0] != "JWT" || got[1] != "scope" {
		t.Fatalf("scoped op must run authn then scope, ran %v", got)
	}
	// an operation without x-scope keeps today's chain — no scope middleware
	if got := run(r["GetUnscoped"]); len(got) != 1 || got[0] != "JWT" {
		t.Fatalf("unscoped op must keep only its auth chain, ran %v", got)
	}
}

func TestChainsScopeRejectionShortCircuits(t *testing.T) {
	useTestProblems(t, func(err error) *fp.Problem {
		if errors.Is(err, errScopeDenied) {
			return fp.New(http.StatusForbidden).WithType("insufficient-role").WithDetail(err.Error())
		}
		return nil
	})
	r := mustChains(t, specWithScope, Primitives{
		JWT:    named("JWT"),
		Scopes: map[string]ScopeChecker{"platform-only": rejectingScope},
	})

	invocations = nil
	rec := httptest.NewRecorder()
	var next http.Handler = http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("scope rejection must stop the chain before the handler")
	})
	chain := r["GetScoped"]
	for i := range slices.Backward(chain) {
		next = chain[i](next)
	}
	next.ServeHTTP(rec, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil))

	if rec.Code != http.StatusForbidden {
		t.Fatalf("scope rejection must write its problem, 403, got %d", rec.Code)
	}
}

func TestChainsScopeOnPublicOperationFails(t *testing.T) {
	_, err := chainsFor(specScopePublic, Primitives{
		JWT:    named("JWT"),
		Scopes: map[string]ScopeChecker{"platform-only": scopeNamed()},
	})
	if err == nil || !strings.Contains(err.Error(), "no security schemes") {
		t.Fatalf("want error for x-scope on a public operation, got %v", err)
	}
}

func TestChainsUnknownScopeFails(t *testing.T) {
	// unknown scope with a non-nil registry
	_, err := chainsFor(specScopeUnknown, Primitives{
		JWT: named("JWT"),
		Scopes: map[string]ScopeChecker{
			"platform-only": scopeNamed(),
		},
	})
	if err == nil || !strings.Contains(err.Error(), "getScopedUnknown") {
		t.Fatalf("want error naming the operation and unknown scope, got %v", err)
	}

	// x-scope with a nil registry — the feature is off, but a spec that
	// declares scopes must not silently pass
	_, err = chainsFor(specScopeUnknown, Primitives{JWT: named("JWT")})
	if err == nil || !strings.Contains(err.Error(), "no scope checker registered") {
		t.Fatalf("want error for x-scope with no registry, got %v", err)
	}
}

func TestChainsNilScopeCheckerFails(t *testing.T) {
	_, err := chainsFor(specWithScope, Primitives{
		JWT:    named("JWT"),
		Scopes: map[string]ScopeChecker{"platform-only": nil},
	})
	if err == nil || !strings.Contains(err.Error(), "no scope checker registered") {
		t.Fatalf("want error for a registered-but-nil scope checker, got %v", err)
	}
}

func cacheSpec(method, security, cache string) string {
	return `
openapi: 3.1.0
info: {title: t, version: "1"}
paths:
  /cached:
    ` + method + `:
      operationId: cached
      security: ` + security + `
      x-cache: ` + cache + `
`
}

func TestChainsSharedCacheMarksTheResponse(t *testing.T) {
	r := mustChains(t, cacheSpec("get", "[]", "shared"), testPrimitives)
	var next http.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	chain := r["Cached"]
	for i := range slices.Backward(chain) {
		next = chain[i](next)
	}
	rec := httptest.NewRecorder()
	httpcache.Middleware(next).ServeHTTP(rec, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/cached", nil))
	if got := rec.Header().Get("Cache-Control"); got != httpcache.SharedClient {
		t.Fatalf("Cache-Control = %q, want %q", got, httpcache.SharedClient)
	}
}

func TestChainsCacheMisuseFails(t *testing.T) {
	for name, spec := range map[string]string{
		"authenticated": cacheSpec("get", "[{bearerAuth: []}]", "shared"),
		"not a GET":     cacheSpec("post", "[]", "shared"),
		"unknown value": cacheSpec("get", "[]", "forever"),
	} {
		_, err := chainsFor(spec, testPrimitives)
		if err == nil || !strings.Contains(err.Error(), "x-cache") {
			t.Errorf("%s: want an x-cache error, got %v", name, err)
		}
	}
}
