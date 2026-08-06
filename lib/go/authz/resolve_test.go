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
)

var testPrimitives = Primitives{
	JWT:    named("JWT"),
	APIKey: named("APIKey"),
	Any:    named("AnyAuth"),
}

func mustResolver(t *testing.T, spec string, primitives Primitives, opts Options) *Resolver {
	t.Helper()
	r, err := NewResolver([]byte(spec), primitives, opts)
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
	r := mustResolver(t, specWithDefault, testPrimitives, Options{})

	if got := run(r.Chains()["GetPublic"]); len(got) != 0 {
		t.Fatalf("getPublic must be public, ran %v", got)
	}
	if got := run(r.Chains()["GetDefaulted"]); len(got) != 1 || got[0] != "JWT" {
		t.Fatalf("getDefaulted must inherit default bearerAuth, ran %v", got)
	}
	if got := run(r.Chains()["GetJWT"]); len(got) != 1 || got[0] != "JWT" {
		t.Fatalf("getJWT must be bearerAuth, ran %v", got)
	}
	if got := run(r.Chains()["GetAPIKey"]); len(got) != 1 || got[0] != "APIKey" {
		t.Fatalf("getAPIKey must be apiKeyAuth, ran %v", got)
	}
	if got := run(r.Chains()["GetAny"]); len(got) != 1 || got[0] != "AnyAuth" {
		t.Fatalf("getAny must derive the OR combination to AnyAuth, ran %v", got)
	}
}

func TestResolverNoDefaultMeansPublic(t *testing.T) {
	r := mustResolver(t, specNoDefault, testPrimitives, Options{})
	if got := run(r.Chains()["GetX"]); len(got) != 0 {
		t.Fatalf("operation without default security must be public, ran %v", got)
	}
}

func TestResolverUnregisteredCombinationFails(t *testing.T) {
	_, err := NewResolver([]byte(specUnknownScheme), testPrimitives, Options{})
	if err == nil {
		t.Fatal("want error for unregistered security combination, got nil")
	}
}

func TestResolverAndBlockFails(t *testing.T) {
	_, err := NewResolver([]byte(specAndBlock), testPrimitives, Options{})
	if err == nil || !strings.Contains(err.Error(), "separate blocks") {
		t.Fatalf("want error for single-block AND security, got %v", err)
	}
}

func TestResolverNilPrimitiveFails(t *testing.T) {
	_, err := NewResolver([]byte(specWithDefault), Primitives{}, Options{})
	if err == nil || !strings.Contains(err.Error(), "primitives.JWT is nil") {
		t.Fatalf("want error for nil primitive, got %v", err)
	}
}

func TestResolverSetupGuard(t *testing.T) {
	r := mustResolver(t, specWithDefault, testPrimitives, Options{
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

func TestResolverUnknownOperationListFails(t *testing.T) {
	_, err := NewResolver([]byte(specWithDefault), testPrimitives, Options{
		SkipSetupGuard: []string{"getNotInSpec"},
	})
	if err == nil || !strings.Contains(err.Error(), "getNotInSpec") {
		t.Fatalf("want error naming the unknown SkipSetupGuard entry, got %v", err)
	}
}
