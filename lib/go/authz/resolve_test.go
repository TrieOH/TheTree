package authz

import (
	"net/http"
	"net/http/httptest"
	"testing"
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
	var next http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	for i := len(chain) - 1; i >= 0; i-- {
		next = chain[i](next)
	}
	next.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
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
)

func mustResolver(t *testing.T, spec string, registry Registry, opts Options) *Resolver {
	t.Helper()
	r, err := NewResolver([]byte(spec), registry, opts)
	if err != nil {
		t.Fatalf("NewResolver: %v", err)
	}
	return r
}

func TestResolverDefaultSecurityAndOverrides(t *testing.T) {
	registry := Registry{
		"bearerAuth":            named("JWT"),
		"apiKeyAuth+bearerAuth": named("AnyAuth"),
	}
	r := mustResolver(t, specWithDefault, registry, Options{})

	if got := run(r.Chains()["getPublic"]); len(got) != 0 {
		t.Fatalf("getPublic must be public, ran %v", got)
	}
	if got := run(r.Chains()["getDefaulted"]); len(got) != 1 || got[0] != "JWT" {
		t.Fatalf("getDefaulted must inherit default bearerAuth, ran %v", got)
	}
	if got := run(r.Chains()["getJWT"]); len(got) != 1 || got[0] != "JWT" {
		t.Fatalf("getJWT must be bearerAuth, ran %v", got)
	}
	if got := run(r.Chains()["getAny"]); len(got) != 1 || got[0] != "AnyAuth" {
		t.Fatalf("getAny must use the OR combination, ran %v", got)
	}
}

func TestResolverNoDefaultMeansPublic(t *testing.T) {
	r := mustResolver(t, specNoDefault, Registry{"bearerAuth": named("JWT")}, Options{})
	if got := run(r.Chains()["getX"]); len(got) != 0 {
		t.Fatalf("operation without default security must be public, ran %v", got)
	}
}

func TestResolverUnregisteredCombinationFails(t *testing.T) {
	_, err := NewResolver([]byte(specUnknownScheme), Registry{"bearerAuth": named("JWT")}, Options{})
	if err == nil {
		t.Fatal("want error for unregistered security combination, got nil")
	}
}

func TestResolverSetupGuard(t *testing.T) {
	registry := Registry{
		"bearerAuth":            named("JWT"),
		"apiKeyAuth+bearerAuth": named("AnyAuth"),
	}
	r := mustResolver(t, specWithDefault, registry, Options{
		SetupGuard:     named("setupGuard"),
		SkipSetupGuard: []string{"getDefaulted"},
	})

	// skipped operation keeps only its auth middleware
	if got := run(r.Chains()["getDefaulted"]); len(got) != 1 || got[0] != "JWT" {
		t.Fatalf("skipped op must keep only its auth chain, ran %v", got)
	}
	// every other operation — public included — gains the guard in front
	if got := run(r.Chains()["getPublic"]); len(got) != 1 || got[0] != "setupGuard" {
		t.Fatalf("public op must still run the setup guard, ran %v", got)
	}
	if got := run(r.Chains()["getJWT"]); len(got) != 2 || got[0] != "setupGuard" || got[1] != "JWT" {
		t.Fatalf("protected op must run guard then auth, ran %v", got)
	}
}
