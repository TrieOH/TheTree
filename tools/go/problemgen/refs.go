package main

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// componentKinds are the components a spec file may declare. Every one is
// merged into the bundle under its own name, so a $ref to it keeps its
// name and only loses its file.
var componentKinds = []string{"schemas", "parameters", "responses", "requestBodies", "headers", "examples"}

// specFile is a file the bundle reads: the spec, a problem file, or a
// feature's spec.yml reached through a $ref.
type specFile struct {
	path  string
	root  *yaml.Node
	paths map[string]bool // the paths it declares that the spec imported
}

// ref is a parsed $ref: the file it points into (absolute; the spec's own
// path for a local ref) and the JSON pointer's segments.
type refTarget struct {
	file    string
	pointer []string
}

func parseRef(value, from string) (refTarget, error) {
	filePart, fragment, ok := strings.Cut(value, "#")
	if !ok {
		return refTarget{}, fmt.Errorf("$ref %s names no #/ pointer", value)
	}
	file := from
	if filePart != "" {
		file = filepath.Join(filepath.Dir(from), filePart)
	}
	fragment, err := url.PathUnescape(fragment)
	if err != nil {
		return refTarget{}, fmt.Errorf("$ref %s: %w", value, err)
	}
	if !strings.HasPrefix(fragment, "/") {
		return refTarget{}, fmt.Errorf("$ref %s: the pointer must start with /", value)
	}
	segs := strings.Split(fragment[1:], "/")
	for i, s := range segs {
		segs[i] = strings.ReplaceAll(strings.ReplaceAll(s, "~1", "/"), "~0", "~")
	}
	return refTarget{file: file, pointer: segs}, nil
}

// component reports whether t points at a component, and which.
func (t refTarget) component() (kind, name string, ok bool) {
	if len(t.pointer) != 3 || t.pointer[0] != "components" {
		return "", "", false
	}
	return t.pointer[1], t.pointer[2], true
}

// importPaths replaces every path item of the spec that is a $ref into
// another file's paths with that item, and loads every file a $ref reaches.
// Each $ref left in the bundle is local: a component keeps its name when
// merged, so the ref only drops its file. The spec reads like any
// multi-file OpenAPI document; only the bundle is one file.
func (b *bundler) importPaths(specPath string) error {
	b.specPath = specPath
	paths := get(b.root, "paths")
	if paths != nil {
		for i := 0; i < len(paths.Content); i += 2 {
			key, item := paths.Content[i].Value, paths.Content[i+1]
			r := get(item, "$ref")
			if r == nil {
				err := b.localize(item, specPath)
				if err != nil {
					return fmt.Errorf("%s: paths %s: %w", specPath, key, err)
				}
				continue
			}
			if len(item.Content) != 2 {
				return fmt.Errorf("%s: paths %s: a $ref path item has nothing beside the $ref", specPath, key)
			}
			imported, err := b.importPath(key, r.Value, specPath)
			if err != nil {
				return fmt.Errorf("%s: paths %s: %w", specPath, key, err)
			}
			paths.Content[i+1] = imported
		}
	}
	err := b.localize(get(b.root, "components"), specPath)
	if err != nil {
		return fmt.Errorf("%s: %w", specPath, err)
	}
	for _, f := range b.files {
		for key, used := range f.paths {
			if !used {
				return fmt.Errorf("%s: paths %s is never imported by %s", f.path, key, specPath)
			}
		}
	}
	return nil
}

// importPath returns the path item a root path's $ref names, its own refs
// made local.
func (b *bundler) importPath(key, value, from string) (*yaml.Node, error) {
	t, err := parseRef(value, from)
	if err != nil {
		return nil, err
	}
	if len(t.pointer) != 2 || t.pointer[0] != "paths" {
		return nil, fmt.Errorf("$ref %s must point at a path item (#/paths/...)", value)
	}
	if t.pointer[1] != key {
		return nil, fmt.Errorf("$ref %s imports %s under another path", value, t.pointer[1])
	}
	f, err := b.loadSpecFile(t.file)
	if err != nil {
		return nil, err
	}
	item := get(get(f.root, "paths"), key)
	if item == nil {
		return nil, fmt.Errorf("$ref %s: %s declares no path %s", value, f.path, key)
	}
	if f.paths[key] {
		return nil, fmt.Errorf("$ref %s: the path is imported twice", value)
	}
	f.paths[key] = true
	err = b.localize(item, f.path)
	if err != nil {
		return nil, fmt.Errorf("%s: paths %s: %w", f.path, key, err)
	}
	return item, nil
}

// loadSpecFile reads a file a $ref reaches, once: its paths wait to be
// imported, its components join the bundle like a problem file's.
func (b *bundler) loadSpecFile(path string) (*specFile, error) {
	path = filepath.Clean(path)
	if f, ok := b.files[path]; ok {
		return f, nil
	}
	if path == filepath.Clean(b.specPath) {
		return nil, errors.New("a path item can't $ref the spec itself")
	}
	err := b.loadFile(path)
	if err != nil {
		return nil, err
	}
	return b.files[path], nil
}

// localize rewrites every $ref under n, read from file, to the bundle's
// local form, loading the files they reach. A ref must name a component
// its file declares.
func (b *bundler) localize(n *yaml.Node, file string) error {
	if n == nil {
		return nil
	}
	if n.Kind == yaml.MappingNode {
		for i := 0; i < len(n.Content); i += 2 {
			if n.Content[i].Value != "$ref" || n.Content[i+1].Kind != yaml.ScalarNode {
				continue
			}
			value := n.Content[i+1].Value
			t, err := parseRef(value, file)
			if err != nil {
				return err
			}
			kind, name, ok := t.component()
			if !ok {
				return fmt.Errorf("$ref %s must point at a component (#/components/<kind>/<name>)", value)
			}
			err = b.checkComponent(t.file, kind, name)
			if err != nil {
				return fmt.Errorf("$ref %s: %w", value, err)
			}
			n.Content[i+1].Value = "#/components/" + kind + "/" + name
		}
	}
	for _, c := range n.Content {
		err := b.localize(c, file)
		if err != nil {
			return err
		}
	}
	return nil
}

// checkComponent loads file if a ref first reaches it there, and checks it
// declares the component.
func (b *bundler) checkComponent(file, kind, name string) error {
	file = filepath.Clean(file)
	var root *yaml.Node
	switch {
	case file == filepath.Clean(b.specPath):
		root = b.root
	default:
		f, err := b.loadSpecFile(file)
		if err != nil {
			return err
		}
		root = f.root
	}
	if get(get(get(root, "components"), kind), name) == nil {
		return fmt.Errorf("%s declares no %s %s", file, kind, name)
	}
	return nil
}
