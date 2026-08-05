package authz

import (
	"context"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
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
)

func mustResolver(t *testing.T, spec string, registry Registry, opts Options) *Resolver {
	t.Helper()
	r, err := NewResolver([]byte(spec), registry, opts)
	if err != nil {
		t.Fatalf("NewResolver: %v", err)
	}
	return r
}

func TestGeneratedOperationID(t *testing.T) {
	cases := []struct{ spec, want string }{
		{"listActors", "ListActors"},
		{"getSetup", "GetSetup"},
		{"", ""},
		{"x", "X"},
	}
	for _, c := range cases {
		if got := GeneratedOperationID(c.spec); got != c.want {
			t.Fatalf("GeneratedOperationID(%q) = %q, want %q", c.spec, got, c.want)
		}
	}
}

func TestResolverDefaultSecurityAndOverrides(t *testing.T) {
	registry := Registry{
		"bearerAuth":            named("JWT"),
		"apiKeyAuth+bearerAuth": named("AnyAuth"),
	}
	r := mustResolver(t, specWithDefault, registry, Options{})

	if got := run(r.Chains()["GetPublic"]); len(got) != 0 {
		t.Fatalf("getPublic must be public, ran %v", got)
	}
	if got := run(r.Chains()["GetDefaulted"]); len(got) != 1 || got[0] != "JWT" {
		t.Fatalf("getDefaulted must inherit default bearerAuth, ran %v", got)
	}
	if got := run(r.Chains()["GetJWT"]); len(got) != 1 || got[0] != "JWT" {
		t.Fatalf("getJWT must be bearerAuth, ran %v", got)
	}
	if got := run(r.Chains()["GetAny"]); len(got) != 1 || got[0] != "AnyAuth" {
		t.Fatalf("getAny must use the OR combination, ran %v", got)
	}
}

func TestResolverNoDefaultMeansPublic(t *testing.T) {
	r := mustResolver(t, specNoDefault, Registry{"bearerAuth": named("JWT")}, Options{})
	if got := run(r.Chains()["GetX"]); len(got) != 0 {
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
	if got := run(r.Chains()["GetDefaulted"]); len(got) != 1 || got[0] != "JWT" {
		t.Fatalf("skipped op must keep only its auth chain, ran %v", got)
	}
	// every other operation — public included — gains the guard in front
	if got := run(r.Chains()["GetPublic"]); len(got) != 1 || got[0] != "setupGuard" {
		t.Fatalf("public op must still run the setup guard, ran %v", got)
	}
	if got := run(r.Chains()["GetJWT"]); len(got) != 2 || got[0] != "setupGuard" || got[1] != "JWT" {
		t.Fatalf("protected op must run guard then auth, ran %v", got)
	}
}

func TestResolverClientOnly(t *testing.T) {
	registry := Registry{
		"bearerAuth":            named("JWT"),
		"apiKeyAuth+bearerAuth": named("AnyAuth"),
	}
	r := mustResolver(t, specWithDefault, registry, Options{
		ClientOnly:    named("ClientOnly"),
		ClientOnlyOps: []string{"getJWT", "getSetup"},
	})

	// client-only op: scheme middleware then the client scope guard
	if got := run(r.Chains()["GetJWT"]); len(got) != 2 || got[0] != "JWT" || got[1] != "ClientOnly" {
		t.Fatalf("client-only protected op must run auth then client guard, ran %v", got)
	}
	// public client-only op: guard alone, still applied
	if got := run(r.Chains()["GetSetup"]); len(got) != 1 || got[0] != "ClientOnly" {
		t.Fatalf("public client-only op must still run the client guard, ran %v", got)
	}
	// op not in the list keeps its chain untouched
	if got := run(r.Chains()["GetAny"]); len(got) != 1 || got[0] != "AnyAuth" {
		t.Fatalf("unlisted op must not gain the client guard, ran %v", got)
	}
}

func TestResolverUnknownOperationListFails(t *testing.T) {
	registry := Registry{"bearerAuth": named("JWT")}

	_, err := NewResolver([]byte(specWithDefault), registry, Options{
		ClientOnlyOps: []string{"getNotInSpec"},
		ClientOnly:    named("ClientOnly"),
	})
	if err == nil || !strings.Contains(err.Error(), "getNotInSpec") {
		t.Fatalf("want error naming the unknown ClientOnlyOps entry, got %v", err)
	}

	_, err = NewResolver([]byte(specWithDefault), registry, Options{
		SkipSetupGuard: []string{"getAlsoNotInSpec"},
	})
	if err == nil || !strings.Contains(err.Error(), "getAlsoNotInSpec") {
		t.Fatalf("want error naming the unknown SkipSetupGuard entry, got %v", err)
	}
}

func TestResolverClientOnlyOpsWithoutMiddlewareFails(t *testing.T) {
	_, err := NewResolver([]byte(specWithDefault), Registry{"bearerAuth": named("JWT")}, Options{
		ClientOnlyOps: []string{"getJWT"},
	})
	if err == nil {
		t.Fatal("want error for ClientOnlyOps without a ClientOnly middleware, got nil")
	}
}
