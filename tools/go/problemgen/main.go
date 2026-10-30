// Command problemgen turns a Backend's declared problems (RFC 9457) into the
// bundled spec and the Backend's problems package (ADR-0019). It runs from
// the Backend's directory, after sqlc and before oapi-codegen:
//
//	go run tools/problemgen api-spec.yml api-spec.gen.yml internal/problems/problems.gen.go
//
// The problems come from two files next to the spec: the shared
// ../problems.yml (api/problems.yml) and the Backend's own problems.yml. A
// problem is a schema with an x-problem block (slug, status, title,
// description) and only its extension members. The bundle merges the files'
// schemas into the spec, expands each problem into a <Name>Problem schema
// and response served as application/problem+json, and builds each
// operation's error responses from its x-problems list (a slug, or a status
// for every problem declared with it), which the bundle keeps as slugs. The
// bundle is the Backend's one spec: codegen, orval and /docs/openapi.yml
// read it.
// The package gets one constructor per problem, with the extensions as
// typed parameters, the docs the harness serves at /problems/{slug}, and
// Mapper.
package main

import (
	"fmt"
	"os"
	"path/filepath"
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
	return os.WriteFile(goOut, src, 0o600) //nolint:gosec // as above
}
