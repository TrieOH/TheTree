package authz

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Registry maps a security-scheme combination to the middleware that
// enforces it. Combinations are the sorted union of the scheme names in an
// operation's security blocks, joined with "+". The empty combination is
// public and never looked up.
//
// A backend registers exactly the combinations its spec uses — e.g.
// "bearerAuth" (JWT), "apiKeyAuth", and "apiKeyAuth+bearerAuth" for the
// OR-style anyAuth middleware. An operation whose combination has no entry
// is an error at construction: the spec and the backend's auth stack
// disagree, and that must fail startup, not production.
type Registry map[string]func(http.Handler) http.Handler

// Options carries per-backend decorations the spec cannot express.
type Options struct {
	// SetupGuard, when set, is prepended to every operation except those
	// named in SkipSetupGuard (spec-form operationIds). identityx uses it
	// to gate everything until /auth/setup has run.
	SetupGuard     func(http.Handler) http.Handler
	SkipSetupGuard []string
	// ClientOnly, when set with ClientOnlyOps, is appended to the chain of
	// every listed operation (spec-form operationIds). identityx uses it to
	// require a platform-level client identity (nil ProjectID) on the
	// platform operations. Both lists are validated against the spec at
	// construction: an unknown operationId fails startup instead of
	// silently dropping a guard.
	ClientOnly    func(http.Handler) http.Handler
	ClientOnlyOps []string
}

// ResolvedOperation is one spec operation with its effective security.
// Schemes is the sorted union of the scheme names across the operation's
// security blocks (or the spec default); empty means public.
type ResolvedOperation struct {
	Method      string
	Path        string
	OperationID string
	Schemes     []string
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
}

// SpecOperations parses the spec and returns every operation with its
// effective security (the spec default applied where the operation declares
// none) and its route (method + path). Both the resolver and the per-backend
// parity tests derive their expectations from this one reading of the spec.
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
			set := make(map[string]bool)
			for _, block := range blocks {
				for scheme := range block {
					set[scheme] = true
				}
			}
			schemes := make([]string, 0, len(set))
			for scheme := range set {
				schemes = append(schemes, scheme)
			}
			sort.Strings(schemes)
			ops = append(ops, ResolvedOperation{
				Method:      m.method,
				Path:        path,
				OperationID: m.op.OperationID,
				Schemes:     schemes,
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

// Resolver derives every operation's middleware chain from the spec's
// security blocks: the spec is the single source of truth for who may call
// what, and with which scheme. Chains are keyed by generated-form
// operationId (see GeneratedOperationID) — the exact value the generated
// strict server passes at dispatch time.
type Resolver struct {
	chains map[string][]func(http.Handler) http.Handler
}

// NewResolver parses the OpenAPI spec and resolves every operation's chain.
// It errors when an operation's security combination has no Registry entry,
// when the spec is not valid YAML, when a named operation in SkipSetupGuard
// or ClientOnlyOps does not exist in the spec, or when ClientOnlyOps is set
// without a ClientOnly middleware — all fail startup, not production.
func NewResolver(spec []byte, registry Registry, opts Options) (*Resolver, error) {
	ops, err := SpecOperations(spec)
	if err != nil {
		return nil, err
	}

	known := make(map[string]bool, len(ops))
	for _, op := range ops {
		known[op.OperationID] = true
	}

	skipSetup := make(map[string]bool, len(opts.SkipSetupGuard))
	for _, id := range opts.SkipSetupGuard {
		if !known[id] {
			return nil, fmt.Errorf("authz resolver: SkipSetupGuard: operation %q not found in spec", id)
		}
		skipSetup[id] = true
	}

	clientOnlyOps := make(map[string]bool, len(opts.ClientOnlyOps))
	for _, id := range opts.ClientOnlyOps {
		if !known[id] {
			return nil, fmt.Errorf("authz resolver: ClientOnlyOps: operation %q not found in spec", id)
		}
		clientOnlyOps[id] = true
	}
	if len(clientOnlyOps) > 0 && opts.ClientOnly == nil {
		return nil, errors.New("authz resolver: ClientOnlyOps set without a ClientOnly middleware")
	}

	r := &Resolver{chains: make(map[string][]func(http.Handler) http.Handler, len(ops))}
	for _, op := range ops {
		var chain []func(http.Handler) http.Handler
		if len(op.Schemes) > 0 {
			key := strings.Join(op.Schemes, "+")
			mw, ok := registry[key]
			if !ok {
				return nil, fmt.Errorf("authz resolver: %s %s: security combination %q not registered", op.Method, op.OperationID, key)
			}
			if mw == nil {
				return nil, fmt.Errorf("authz resolver: %s %s: middleware for combination %q is nil", op.Method, op.OperationID, key)
			}
			chain = []func(http.Handler) http.Handler{mw}
		}
		if opts.SetupGuard != nil && !skipSetup[op.OperationID] {
			chain = append([]func(http.Handler) http.Handler{opts.SetupGuard}, chain...)
		}
		if clientOnlyOps[op.OperationID] {
			chain = append(chain, opts.ClientOnly)
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
