package authz

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/MintzyG/fun"
	"gopkg.in/yaml.v3"
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
// composition, e.g. fun's AnyAuth) when both are listed. There is no
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
	Scopes map[string]ScopeChecker          // enforces per-operation x-scope
}

// Options carries per-backend decorations the spec cannot express.
type Options struct {
	// SetupGuard, when set, is prepended to every operation except those
	// named in SkipSetupGuard (spec-form operationIds). identityx uses it
	// to gate everything until /auth/setup has run.
	SetupGuard     func(http.Handler) http.Handler
	SkipSetupGuard []string
}

// ResolvedOperation is one spec operation with its effective security.
// Blocks are the operation's security blocks (the spec default applied
// where the operation declares none): a list of blocks, each a list of
// scheme names. Schemes is their flattened sorted union; empty means
// public. Scope is the operation's x-scope value, empty when it declares
// none (no scope restriction).
type ResolvedOperation struct {
	Method      string
	Path        string
	OperationID string
	Blocks      [][]string
	Schemes     []string
	Scope       string
}

type openAPISpec struct {
	Security []map[string][]string `yaml:"security"`
	Paths    map[string]pathItem   `yaml:"paths"`
}

type pathItem struct {
	Get     *operation `yaml:"get"`
	Put     *operation `yaml:"put"`
	Post    *operation `yaml:"post"`
	Patch   *operation `yaml:"patch"`
	Delete  *operation `yaml:"delete"`
	Options *operation `yaml:"options"`
	Head    *operation `yaml:"head"`
}

type operation struct {
	OperationID string                `yaml:"operationId"`
	Security    []map[string][]string `yaml:"security"`
	Scope       *string               `yaml:"x-scope"`
}

// SpecOperations parses the spec and returns every operation with its
// effective security (the spec default applied where the operation declares
// none), its x-scope (empty when undeclared), and its route (method + path).
// Both the resolver and the per-backend parity tests derive their
// expectations from this one reading of the spec.
func SpecOperations(spec []byte) ([]ResolvedOperation, error) {
	var doc openAPISpec
	err := yaml.Unmarshal(spec, &doc)
	if err != nil {
		return nil, fmt.Errorf("authz: parse spec: %w", err)
	}
	var ops []ResolvedOperation
	for path, item := range doc.Paths {
		for _, m := range []struct {
			method string
			op     *operation
		}{{"GET", item.Get}, {"PUT", item.Put}, {"POST", item.Post}, {"PATCH", item.Patch}, {"DELETE", item.Delete}, {"OPTIONS", item.Options}, {"HEAD", item.Head}} {
			if m.op == nil || m.op.OperationID == "" {
				continue
			}
			blocks := m.op.Security
			if blocks == nil {
				blocks = doc.Security
			}
			var names [][]string
			set := make(map[string]bool)
			for _, block := range blocks {
				schemes := make([]string, 0, len(block))
				for scheme := range block {
					schemes = append(schemes, scheme)
				}
				sort.Strings(schemes)
				names = append(names, schemes)
				for _, s := range schemes {
					set[s] = true
				}
			}
			union := make([]string, 0, len(set))
			for s := range set {
				union = append(union, s)
			}
			sort.Strings(union)
			ops = append(ops, ResolvedOperation{
				Method:      m.method,
				Path:        path,
				OperationID: m.op.OperationID,
				Blocks:      names,
				Schemes:     union,
				Scope:       scopeOf(m.op.Scope),
			})
		}
	}
	return ops, nil
}

// GeneratedOperationID returns the operationID form the oapi-codegen strict
// server passes to strict middlewares (its method name on the generated
// interface): spec-form "listActors" becomes "ListActors". Chains are keyed
// by this form so the dispatch can do an exact lookup with no runtime
// string surgery — the inverse of the old per-backend casing hack, whose
// failure mode was a miss being treated as "public".
func GeneratedOperationID(specID string) string {
	if specID == "" {
		return specID
	}
	return strings.ToUpper(specID[:1]) + specID[1:]
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

// scopeOf converts the declared x-scope to its string form: absent (nil)
// is the empty scope, so operations without x-scope carry no restriction.
func scopeOf(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// scopeMiddleware wraps a scope checker into an HTTP middleware that sends
// the checker's error on rejection and otherwise continues the chain. It is
// the only place a scope predicate meets HTTP; the predicate itself stays
// context-only and testable.
func scopeMiddleware(checker ScopeChecker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := checker(r.Context()); err != nil {
				fun.Error(err).Send(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// Resolver derives every operation's middleware chain from the spec's
// security blocks: the spec is the single source of truth for who may call
// what, and with which scheme. Chains are keyed by generated-form
// operationId (see GeneratedOperationID) — the exact value the generated
// strict server passes at dispatch time.
type Resolver struct {
	chains map[string][]func(http.Handler) http.Handler
}

// NewResolver parses the OpenAPI spec and resolves every operation's chain.
// It errors when an operation's security combination has no primitive (or a
// nil one), when a single security block lists multiple schemes (AND
// semantics — the resolver composes OR across blocks; write the spec's
// alternatives as separate blocks), when an operation declares an x-scope
// with no registered checker in primitives.Scopes, when an operation
// declares an x-scope but no security schemes (a scope reads the identity
// authn wrote — it implies authentication), when the spec is not valid YAML,
// or when a named operation in SkipSetupGuard does not exist in the spec —
// all fail startup, not production.
func NewResolver(spec []byte, primitives Primitives, opts Options) (*Resolver, error) {
	ops, err := SpecOperations(spec)
	if err != nil {
		return nil, err
	}

	known := make(map[string]bool, len(ops))
	for _, op := range ops {
		known[op.OperationID] = true
		for _, block := range op.Blocks {
			if len(block) > 1 {
				return nil, fmt.Errorf("authz resolver: %s %s: security block lists multiple schemes %v (AND) — express OR as separate blocks", op.Method, op.OperationID, block)
			}
		}
		if op.Scope != "" {
			if len(op.Schemes) == 0 {
				return nil, fmt.Errorf("authz resolver: %s %s: x-scope %q on an operation with no security schemes (a scope implies authentication)", op.Method, op.OperationID, op.Scope)
			}
			if primitives.Scopes == nil || primitives.Scopes[op.Scope] == nil {
				return nil, fmt.Errorf("authz resolver: %s %s: no scope checker registered for x-scope %q", op.Method, op.OperationID, op.Scope)
			}
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
			chain = append(chain, scopeMiddleware(primitives.Scopes[op.Scope]))
		}
		if opts.SetupGuard != nil && !skipSetup[op.OperationID] {
			chain = append([]func(http.Handler) http.Handler{opts.SetupGuard}, chain...)
		}
		r.chains[GeneratedOperationID(op.OperationID)] = chain
	}
	return r, nil
}

// Chains returns every resolved chain, keyed by generated-form operationId
// (the form the generated strict server passes to its middlewares). Public
// operations carry an empty chain; operations absent from the map were
// absent from the spec.
func (r *Resolver) Chains() map[string][]func(http.Handler) http.Handler {
	return r.chains
}
