package authz

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
// directly without HTTP, and the resolver wraps them into middleware.
type ScopeChecker func(context.Context) error

// Primitives holds the auth middlewares a backend provides. The resolver
// derives each operation's chain from the spec's security blocks: the union
// of the scheme names an operation lists selects the primitive that enforces
// it — JWT for "bearerAuth", APIKey for "apiKeyAuth", Any (the OR
// composition, httpserver.Authenticators' Any) when both are listed. There is no
// per-combination registry to keep in sync with the spec; the spec is the
// single source of truth for who may call what, and with which scheme. An
// operation whose combination has no primitive is an error at construction:
// the spec and the backend's auth stack disagree, and that must fail
// startup, not production.
//
// Scopes is the registry of scope checkers, keyed by the x-scope value the
// spec declares per operation. A nil map disables scope support entirely —
// backends that declare no x-scope never need one. Every x-scope an
// operation declares must name a registered checker; a miss fails startup
// alongside the scheme checks, never production.
type Primitives struct {
	JWT    func(http.Handler) http.Handler // enforces scheme "bearerAuth"
	APIKey func(http.Handler) http.Handler // enforces scheme "apiKeyAuth"
	Any    func(http.Handler) http.Handler // enforces any combination of the two (OR)
	Scopes map[string]ScopeChecker         // enforces per-operation x-scope
}

// Options carries per-backend decorations the spec cannot express.
type Options struct {
	// SetupGuard, when set, is prepended to every operation except those
	// named in SkipSetupGuard (spec-form operationIds). identityx uses it
	// to gate everything until /auth/setup has run.
	SetupGuard     func(http.Handler) http.Handler
	SkipSetupGuard []string
	// WriteError writes a scope checker's rejection. It is required when
	// the spec has an x-scope; lib/backend passes lib/httpserver's writer,
	// which authz cannot import (lib/httpserver imports authz).
	WriteError func(http.ResponseWriter, *http.Request, error)
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
func scopeMiddleware(checker ScopeChecker, writeError func(http.ResponseWriter, *http.Request, error)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := checker(r.Context())
			if err != nil {
				writeError(w, r, err)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// Resolver derives every operation's middleware chain from the spec's
// security blocks: the spec is the single source of truth for who may call
// what, and with which scheme. Chains are keyed by generated-form
// operationId (see openapispec.GeneratedID) — the exact value the generated
// strict server passes at dispatch time.
type Resolver struct {
	chains map[string][]func(http.Handler) http.Handler
}

// NewResolver parses the OpenAPI spec and resolves every operation's chain.
// It errors when an operation's security combination has no primitive (or a
// nil one), when a single security block lists multiple schemes (AND
// semantics — the resolver composes OR across blocks; write the spec's
// alternatives as separate blocks), when an operation declares an x-scope
// with no registered checker in primitives.Scopes or no Options.WriteError,
// when an operation
// declares an x-scope but no security schemes (a scope reads the identity
// authn wrote — it implies authentication), when the spec is not valid YAML,
// when an operation declares an x-cache other than "shared", or on anything
// but a public GET (a shared response is the same for every caller), or
// when a named operation in SkipSetupGuard does not exist in the spec —
// all fail startup, not production.
func NewResolver(spec []byte, primitives Primitives, opts Options) (*Resolver, error) {
	ops, err := openapispec.Operations(spec)
	if err != nil {
		return nil, err
	}

	known := make(map[string]bool, len(ops))
	for _, op := range ops {
		known[op.OperationID] = true
		err = validateOperation(op, primitives, opts)
		if err != nil {
			return nil, err
		}
	}

	skipSetup := make(map[string]bool, len(opts.SkipSetupGuard))
	for _, id := range opts.SkipSetupGuard {
		if !known[id] {
			return nil, fmt.Errorf("authz resolver: SkipSetupGuard: operation %q not found in spec", id)
		}
		skipSetup[id] = true
	}

	r := &Resolver{chains: make(map[string][]func(http.Handler) http.Handler, len(ops))}
	for _, op := range ops {
		var chain []func(http.Handler) http.Handler
		if len(op.Schemes) > 0 {
			mw, err := primitives.forCombination(strings.Join(op.Schemes, "+"))
			if err != nil {
				return nil, fmt.Errorf("authz resolver: %s %s: %w", op.Method, op.OperationID, err)
			}
			chain = []func(http.Handler) http.Handler{mw}
		}
		if op.Scope != "" {
			// A scope reads the identity authn wrote, so it runs after the
			// authn primitive; the setup guard stays outermost.
			chain = append(chain, scopeMiddleware(primitives.Scopes[op.Scope], opts.WriteError))
		}
		if op.Cache == cacheShared {
			chain = append(chain, httpcache.Shared)
		}
		if opts.SetupGuard != nil && !skipSetup[op.OperationID] {
			chain = append([]func(http.Handler) http.Handler{opts.SetupGuard}, chain...)
		}
		r.chains[openapispec.GeneratedID(op.OperationID)] = chain
	}
	return r, nil
}

// validateOperation checks one operation's declarations against what the
// backend registered: OR-only security blocks, x-cache only on a public GET,
// and x-scope only with authentication, a registered checker and a writer.
func validateOperation(op openapispec.Operation, primitives Primitives, opts Options) error {
	for _, block := range op.Blocks {
		if len(block) > 1 {
			return fmt.Errorf("authz resolver: %s %s: security block lists multiple schemes %v (AND) — express OR as separate blocks", op.Method, op.OperationID, block)
		}
	}
	if op.Cache != "" {
		if op.Cache != cacheShared {
			return fmt.Errorf("authz resolver: %s %s: unknown x-cache %q (want %q)", op.Method, op.OperationID, op.Cache, cacheShared)
		}
		if op.Method != http.MethodGet || len(op.Schemes) > 0 {
			return fmt.Errorf("authz resolver: %s %s: x-cache %q on an operation that is not a public GET", op.Method, op.OperationID, op.Cache)
		}
	}
	if op.Scope == "" {
		return nil
	}
	if len(op.Schemes) == 0 {
		return fmt.Errorf("authz resolver: %s %s: x-scope %q on an operation with no security schemes (a scope implies authentication)", op.Method, op.OperationID, op.Scope)
	}
	if primitives.Scopes == nil || primitives.Scopes[op.Scope] == nil {
		return fmt.Errorf("authz resolver: %s %s: no scope checker registered for x-scope %q", op.Method, op.OperationID, op.Scope)
	}
	if opts.WriteError == nil {
		return fmt.Errorf("authz resolver: %s %s: x-scope %q needs Options.WriteError", op.Method, op.OperationID, op.Scope)
	}
	return nil
}

// Chains returns every resolved chain, keyed by generated-form operationId
// (the form the generated strict server passes to its middlewares). Public
// operations carry an empty chain; operations absent from the map were
// absent from the spec.
func (r *Resolver) Chains() map[string][]func(http.Handler) http.Handler {
	return r.chains
}
