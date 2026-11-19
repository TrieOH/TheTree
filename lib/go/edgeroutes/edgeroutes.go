// Package edgeroutes renders the route list an API's Cloudflare Worker
// checks before it starts the container: the first path segment of every
// operation in the service's OpenAPI spec, plus the paths httpserver serves
// outside the spec. A request whose first segment is not on the list cannot
// reach a handler, so the Worker answers 404 without waking the container.
//
// The list is committed as cloudflare/src/routes.ts, next to the Worker
// (each API's deploy job has no Go to generate it). tools/problemgen writes
// it on every `just gen`, and each service's test fails when the committed
// list no longer matches the spec.
package edgeroutes

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// FileName is where each service keeps the rendered list, relative to the
// service module.
const FileName = "cloudflare/src/routes.ts"

// served are the first segments httpserver.NewRouter adds to every backend,
// whether or not its spec declares them (pinned by this package's tests).
var served = []string{"docs", "health", "problems"}

// Segments returns the sorted first path segments of spec's operations,
// together with the segments httpserver serves outside the spec.
func Segments(spec []byte) ([]string, error) {
	var doc struct {
		Paths map[string]any `yaml:"paths"`
	}
	err := yaml.Unmarshal(spec, &doc)
	if err != nil {
		return nil, fmt.Errorf("edgeroutes: parse spec: %w", err)
	}
	if len(doc.Paths) == 0 {
		return nil, errors.New("edgeroutes: spec declares no paths")
	}
	set := map[string]bool{}
	for _, s := range served {
		set[s] = true
	}
	for path := range doc.Paths {
		first, _, _ := strings.Cut(strings.TrimPrefix(path, "/"), "/")
		if first == "" || strings.Contains(first, "{") {
			return nil, fmt.Errorf("edgeroutes: path %q has no fixed first segment", path)
		}
		set[first] = true
	}
	out := make([]string, 0, len(set))
	for s := range set {
		out = append(out, s)
	}
	slices.Sort(out)
	return out, nil
}

// Render returns the TypeScript module the Worker imports.
func Render(spec []byte) ([]byte, error) {
	segments, err := Segments(spec)
	if err != nil {
		return nil, err
	}
	var b strings.Builder
	b.WriteString("// Code generated from api-spec.yml by `just gen`. DO NOT EDIT.\n")
	b.WriteString("//\n")
	b.WriteString("// The first path segments this API serves. The Worker answers any other\n")
	b.WriteString("// path with 404 without starting the container.\n")
	b.WriteString("export const ROUTES: ReadonlySet<string> = new Set([\n")
	for _, s := range segments {
		b.WriteString("\t" + strconv.Quote(s) + ",\n")
	}
	b.WriteString("]);\n")
	return []byte(b.String()), nil
}

// Write renders spec's list into dir (the service module).
func Write(spec []byte, dir string) error {
	want, err := Render(spec)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Clean(filepath.Join(dir, FileName)), want, 0o600)
}

// Verify compares the committed list in dir (the service module) with spec.
func Verify(spec []byte, dir string) error {
	want, err := Render(spec)
	if err != nil {
		return err
	}
	got, err := os.ReadFile(filepath.Clean(filepath.Join(dir, FileName)))
	if err != nil {
		return fmt.Errorf("edgeroutes: %w (run `just gen`)", err)
	}
	if !bytes.Equal(got, want) {
		return fmt.Errorf("edgeroutes: %s is out of date with api-spec.yml (run `just gen`)", FileName)
	}
	return nil
}
