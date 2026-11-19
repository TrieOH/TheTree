// Package openapispec reads what the Backends need from a bundled OpenAPI
// spec, per operation: its security, x-scope, x-problems and x-cache. It is
// the one reading of the spec that the auth chains, the problem check and
// the edge cache marking share.
package openapispec

import (
	"fmt"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Operation is one spec operation with its effective security.
// Blocks are the operation's security blocks (the spec default applied
// where the operation declares none): a list of blocks, each a list of
// scheme names. Schemes is their flattened sorted union; empty means
// public. Scope is the operation's x-scope value, empty when it declares
// none (no scope restriction). Problems is the bundled spec's x-problems,
// the slugs the operation may answer with (ADR-0019). Cache is the
// operation's x-cache value, empty when it declares none (the response is
// private); "shared" lets the edge keep it until a write purges a tag it
// carries (ADR-0020).
type Operation struct {
	Method      string
	Path        string
	OperationID string
	Blocks      [][]string
	Schemes     []string
	Scope       string
	Problems    []string
	Cache       string
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
	Problems    []string              `yaml:"x-problems"`
	Cache       string                `yaml:"x-cache"`
}

// Operations parses the spec and returns every operation with its
// effective security (the spec default applied where the operation declares
// none), its x-scope (empty when undeclared), and its route (method + path).
// lib/backend reads it once at boot, for the auth chains and the problem
// check (both lib/httpserver).
func Operations(spec []byte) ([]Operation, error) {
	var doc openAPISpec
	err := yaml.Unmarshal(spec, &doc)
	if err != nil {
		return nil, fmt.Errorf("openapispec: parse spec: %w", err)
	}
	var ops []Operation
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
			ops = append(ops, Operation{
				Method:      m.method,
				Path:        path,
				OperationID: m.op.OperationID,
				Blocks:      names,
				Schemes:     union,
				Scope:       scopeOf(m.op.Scope),
				Problems:    m.op.Problems,
				Cache:       m.op.Cache,
			})
		}
	}
	return ops, nil
}

// GeneratedID returns the operationID form the oapi-codegen strict
// server passes to strict middlewares (its method name on the generated
// interface): spec-form "listActors" becomes "ListActors". Chains are keyed
// by this form so the dispatch can do an exact lookup with no runtime
// string surgery — the inverse of the old per-backend casing hack, whose
// failure mode was a miss being treated as "public".
func GeneratedID(specID string) string {
	if specID == "" {
		return specID
	}
	return strings.ToUpper(specID[:1]) + specID[1:]
}

// scopeOf converts the declared x-scope to its string form: absent (nil)
// is the empty scope, so operations without x-scope carry no restriction.
func scopeOf(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
