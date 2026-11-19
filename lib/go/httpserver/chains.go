package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"lib/httpcache"
	"lib/openapispec"
)

// ScopeChecker is a predicate answering "is the caller allowed by this
// scope". It reads the identity the auth middleware wrote to the context
// and returns an error to reject; nil passes. Scope checkers are plain
// predicates — they never write to the context — so they are testable
// directly without HTTP, and Chains wraps them into middleware.
type ScopeChecker func(context.Context) error

// Primitives holds the auth middlewares a backend provides. Chains derives
// each operation's chain from the spec's security blocks: the union of the
// scheme names an operation lists selects the primitive that enforces it —
// JWT for "bearerAuth", APIKey for "apiKeyAuth", Any (the OR composition,
// Authenticators' Any) when both are listed. There is no per-combination
// registry to keep in sync with the spec; the spec is the single source of
// truth for who may call what, and with which scheme. An operation whose
// combination has no primitive is an error at boot: the spec and the
// backend's auth stack disagree, and that must fail startup, not production.
//
// Scopes is the registry of scope checkers, keyed by the x-scope value the
// spec declares per operation. Backends that declare no x-scope leave it
// nil. Every x-scope an operation declares must name a registered checker;
// a miss fails startup alongside the scheme checks, never production.
type Primitives struct {
	JWT    func(http.Handler) http.Handler // enforces scheme "bearerAuth"
	APIKey func(http.Handler) http.Handler // enforces scheme "apiKeyAuth"
	Any    func(http.Handler) http.Handler // enforces any combination of the two (OR)
	Scopes map[string]ScopeChecker         // enforces per-operation x-scope
}

// forCombination returns the primitive enforcing a scheme combination (the
// sorted union of scheme names joined with "+"), or an error when the
// backend registered no primitive for it or the primitive is nil. The empty
// combination is public and never looked up.
func (p Primitives) forCombination(combination string) (func(http.Handler) http.Handler, error) {
	switch combination {
	case "bearerAuth":
		if p.JWT == nil {
			return nil, errors.New("primitives.JWT is nil")
		}
		return p.JWT, nil
	case "apiKeyAuth":
		if p.APIKey == nil {
			return nil, errors.New("primitives.APIKey is nil")
		}
		return p.APIKey, nil
	case "apiKeyAuth+bearerAuth":
		if p.Any == nil {
			return nil, errors.New("primitives.Any is nil")
		}
		return p.Any, nil
	default:
		return nil, fmt.Errorf("no primitive for security combination %q", combination)
	}
}

// cacheShared is the one x-cache value: the edge may keep the response.
const cacheShared = "shared"

// scopeMiddleware wraps a scope checker into an HTTP middleware that writes
// the checker's error on rejection and otherwise continues the chain. It is
// the only place a scope predicate meets HTTP; the predicate itself stays
// context-only and testable.
func scopeMiddleware(checker ScopeChecker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := checker(r.Context())
			if err != nil {
				WriteError(w, r, err)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// Chains derives every operation's middleware chain from the spec's
// security, x-scope and x-cache, keyed by generated-form operationId (see
// openapispec.GeneratedID) — the exact value the generated strict server
// passes at dispatch time, and what AuthDispatch looks up. Public operations
// carry an empty chain.
//
// It errors when an operation's security combination has no primitive (or
// a nil one), when a single security block lists multiple schemes (AND
// semantics — chains compose OR across blocks; write the spec's
// alternatives as separate blocks), when an operation declares an x-scope
// with no registered checker or no security schemes (a scope reads the
// identity authn wrote — it implies authentication), and when an operation
// declares an x-cache other than "shared", or one on anything but a public
// GET (a shared response is the same for every caller). All of these fail
// startup, not production.
func Chains(ops []openapispec.Operation, primitives Primitives) (map[string][]func(http.Handler) http.Handler, error) {
	chains := make(map[string][]func(http.Handler) http.Handler, len(ops))
	for _, op := range ops {
		err := validateOperation(op, primitives)
		if err != nil {
			return nil, err
		}
		var chain []func(http.Handler) http.Handler
		if len(op.Schemes) > 0 {
			mw, err := primitives.forCombination(strings.Join(op.Schemes, "+"))
			if err != nil {
				return nil, fmt.Errorf("auth chains: %s %s: %w", op.Method, op.OperationID, err)
			}
			chain = []func(http.Handler) http.Handler{mw}
		}
		if op.Scope != "" {
			// A scope reads the identity authn wrote, so it runs after the
			// authn primitive.
			chain = append(chain, scopeMiddleware(primitives.Scopes[op.Scope]))
		}
		if op.Cache == cacheShared {
			chain = append(chain, httpcache.Shared)
		}
		chains[openapispec.GeneratedID(op.OperationID)] = chain
	}
	return chains, nil
}

// validateOperation checks one operation's declarations against what the
// backend registered: OR-only security blocks, x-cache only on a public GET,
// and x-scope only with authentication and a registered checker.
func validateOperation(op openapispec.Operation, primitives Primitives) error {
	for _, block := range op.Blocks {
		if len(block) > 1 {
			return fmt.Errorf("auth chains: %s %s: security block lists multiple schemes %v (AND) — express OR as separate blocks", op.Method, op.OperationID, block)
		}
	}
	if op.Cache != "" {
		if op.Cache != cacheShared {
			return fmt.Errorf("auth chains: %s %s: unknown x-cache %q (want %q)", op.Method, op.OperationID, op.Cache, cacheShared)
		}
		if op.Method != http.MethodGet || len(op.Schemes) > 0 {
			return fmt.Errorf("auth chains: %s %s: x-cache %q on an operation that is not a public GET", op.Method, op.OperationID, op.Cache)
		}
	}
	if op.Scope == "" {
		return nil
	}
	if len(op.Schemes) == 0 {
		return fmt.Errorf("auth chains: %s %s: x-scope %q on an operation with no security schemes (a scope implies authentication)", op.Method, op.OperationID, op.Scope)
	}
	if primitives.Scopes[op.Scope] == nil {
		return fmt.Errorf("auth chains: %s %s: no scope checker registered for x-scope %q", op.Method, op.OperationID, op.Scope)
	}
	return nil
}
