// Command problemgen turns a Backend's declared problems (RFC 9457) into the
// bundled spec and the Backend's problems package (ADR-0019). It runs from
// the Backend's directory, after sqlc and before oapi-codegen:
//
//	go run tools/problemgen api-spec.yml api-spec.gen.yml internal/problems/problems.gen.go
//
// The spec is split by feature: each path item is a $ref into its feature's
// internal/features/<name>/spec.yml, which declares the path, the
// components only that feature uses and the problems only it answers. Refs
// are ordinary relative OpenAPI refs (a feature reaches the spec's shared
// components as ../../../api-spec.yml#/components/...). The bundle inlines
// the paths, merges every file's components under their own names and
// leaves only local refs; a feature path the spec never imports fails.
//
// The problems come from the shared ../problems.yml (api/problems.yml), the
// Backend's own problems.yml and every feature's spec.yml. A
// problem is a schema with an x-problem block (slug, status, title,
// description) and only its extension members. The bundle merges the files'
// schemas into the spec, expands each problem into a <Name>Problem schema
// and response served as application/problem+json, and builds each
// operation's error responses from its x-problems list (a slug, or a status
// for every problem declared with it), which the bundle keeps as slugs. The
// bundle is the Backend's one spec: codegen, orval and /docs/openapi.yml
// read it. Its info.version comes from the release tags (specVersion).
// The package gets one constructor per problem, with the extensions as
// typed parameters, the docs the harness serves at /problems/{slug}, and
// Mapper.
//
// It also writes the route list the Backend's Worker checks before it wakes
// the container (cloudflare/src/routes.ts, committed; lib/edgeroutes).
//
// A problem may name its Go error and the error's text: x-problem's error
// and detail, the text clients read. Without members the error is a
// sentinel (ErrAddSelf); with members it is a typed error
// (IntentStatusError), a struct with a field per member, and its detail a
// template naming members in braces ({intent_status}). x-shared-errors
// names errors answered by a shared problem the same way. Every one is
// generated into the Backend's errs package (errs/errs.gen.go, next to the
// spec) and Mapper answers it with its problem, so the Backend's own mapper
// keeps only the errors the spec can't express. A detail that only repeats
// the title fails: the detail says what went wrong and what to do.
package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"lib/edgeroutes"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "usage: problemgen <api-spec.yml> <bundled spec out> <problems package file out>")
		os.Exit(2)
	}
	err := run(os.Args[1], os.Args[2], os.Args[3])
	if err != nil {
		fmt.Fprintln(os.Stderr, "problemgen:", err)
		os.Exit(1)
	}
}

func run(specPath, specOut, goOut string) error {
	b, err := load(specPath)
	if err != nil {
		return err
	}
	module, err := modulePath(filepath.Dir(specPath))
	if err != nil {
		return err
	}
	b.module = module
	if v := specVersion(filepath.Dir(specPath)); v != "" {
		set(ensure(b.root, "info"), "version", node(v))
	}
	spec, err := b.bundle()
	if err != nil {
		return err
	}
	src, err := b.goSource()
	if err != nil {
		return err
	}
	err = os.WriteFile(specOut, spec, 0o600) //nolint:gosec // codegen step run by just gen and the Dockerfiles; the path is the build's own
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(goOut), 0o750) //nolint:gosec // as above
	if err != nil {
		return err
	}
	err = os.WriteFile(goOut, src, 0o600) //nolint:gosec // as above
	if err != nil {
		return err
	}
	errsSrc, err := b.errsSource()
	if err != nil {
		return err
	}
	err = writeErrs(filepath.Join(filepath.Dir(specPath), errsFile), errsSrc)
	if err != nil {
		return err
	}
	return writeRoutes(specPath)
}

// specVersion is the release the bundle's info.version names: the tags are
// its one source, so api-spec.yml only holds a placeholder. SPEC_VERSION
// wins (the Dockerfiles get it from the release tag, the front pipeline
// from the remote's tags); without it, the Backend's latest <svc>/v* tag in
// the local repo; with neither (no git), "" keeps the placeholder.
func specVersion(dir string) string {
	if v := os.Getenv("SPEC_VERSION"); v != "" {
		return strings.TrimPrefix(v, "v")
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return ""
	}
	svc := filepath.Base(abs)
	//nolint:gosec // codegen step; svc is the Backend's own directory name
	out, err := exec.CommandContext(context.Background(), "git", "-C", abs, "tag", "-l", svc+"/v*", "--sort=-v:refname").Output()
	if err != nil {
		return ""
	}
	latest, _, _ := strings.Cut(string(out), "\n")
	return strings.TrimPrefix(latest, svc+"/v")
}

// writeErrs writes the generated errs file, or removes it when the Backend
// declares no errors.
func writeErrs(path string, src []byte) error {
	if src == nil {
		err := os.Remove(path) //nolint:gosec // codegen step; the path is the build's own
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return err
	}
	err := os.MkdirAll(filepath.Dir(path), 0o750) //nolint:gosec // as above
	if err != nil {
		return err
	}
	return os.WriteFile(path, src, 0o600) //nolint:gosec // as above
}

// writeRoutes renders the route list the Backend's Worker checks before it
// wakes the container (lib/edgeroutes), when the Backend has a Worker.
func writeRoutes(specPath string) error {
	dir := filepath.Dir(specPath)
	_, err := os.Stat(filepath.Join(dir, filepath.Dir(edgeroutes.FileName))) //nolint:gosec // codegen step; the path is the build's own
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	spec, err := os.ReadFile(specPath) //nolint:gosec // the Backend's own spec
	if err != nil {
		return err
	}
	return edgeroutes.Write(spec, dir)
}

// modulePath is the Go module the spec's Backend is, from its go.mod: the
// generated mapper imports its errs package.
func modulePath(dir string) (string, error) {
	raw, err := os.ReadFile(filepath.Join(dir, "go.mod")) //nolint:gosec // the Backend's own go.mod
	if errors.Is(err, fs.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	for line := range strings.SplitSeq(string(raw), "\n") {
		if rest, ok := strings.CutPrefix(strings.TrimSpace(line), "module "); ok {
			return strings.TrimSpace(rest), nil
		}
	}
	return "", errors.New("go.mod declares no module")
}
