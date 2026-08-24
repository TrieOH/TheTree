// A generated module for Thetree functions

package main

import (
	"context"
	"dagger/thetree/internal/dagger"
	"fmt"
	"strings"
	"sync"

)

type Thetree struct{}

var services = []string{"identityx", "informd", "payssage", "univents"}

const (
	goVersion           = "1.26.4"
	sqlcVersion         = "1.31.1"
	gotestsumVersion    = "1.13.0"
	golangciLintVersion = "2.12.2"
	trivyVersion        = "0.74.0" // keep in sync with trivy-scan.yml filter paths

	// oapiCodegenVersion must stay in sync with the binary baked into the
	// go-tools image (git.trieoh.com/trieoh/go-tools:3), which the service
	// Dockerfiles copy for their builds.
	oapiCodegenVersion = "2.8.0"
)

// goBase returns a bare Go toolchain container with the workspace build
// caches mounted but NO source. Tool-install layers built on top of it
// stay cacheable across commits: dagger invalidates an exec layer whenever
// a mounted source directory changes, so installing tools AFTER mounting
// the source (as before) re-downloaded sqlc / oapi-codegen / golangci-lint
// / gotestsum on every commit.
func (m *Thetree) goBase() *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("golang:%s-bookworm", goVersion)).
		WithMountedCache("/root/go/pkg/mod", dag.CacheVolume("go-mod")).
		WithMountedCache("/root/.cache/go-build", dag.CacheVolume("go-build"))
}

// withSource mounts the workspace source onto a tool-ready container.
// Do this LAST, after every tool install, so the install layers keep their
// cache across source changes.
func (m *Thetree) withSource(c *dagger.Container, source *dagger.Directory, service string) *dagger.Container {
	return c.
		WithDirectory("/workspace", m.scopedSource(source, service)).
		WithWorkdir("/workspace").
		WithEnvVariable("GOWORK", "/workspace/go.work")
}

// baseGo returns a container with the full repo source mounted.
func (m *Thetree) baseGo(source *dagger.Directory) *dagger.Container {
	return m.withSource(m.goBase(), source, "")
}

// baseGoScoped returns a container with only the workspace modules needed
// by the given service. Pass an empty service string for the full source.
func (m *Thetree) baseGoScoped(source *dagger.Directory, service string) *dagger.Container {
	return m.withSource(m.goBase(), source, service)
}

// scopedSource extracts only the Go workspace files — all api/, lib/go/, sdk/go/,
// .dagger/, go.work, and go.work.sum. Skips front/, lib/ts/, sdk/ts/, .git/, etc.
// Pass an empty service string to get the full source directory.
func (m *Thetree) scopedSource(source *dagger.Directory, service string) *dagger.Directory {
	if service == "" {
		return source
	}
	return dag.Directory().
		WithFile("go.work", source.File("go.work")).
		WithFile("go.work.sum", source.File("go.work.sum")).
		WithDirectory("api", source.Directory("api")).
		WithDirectory("lib/go", source.Directory("lib/go")).
		WithDirectory("sdk/go", source.Directory("sdk/go")).
		WithDirectory(".dagger", source.Directory(".dagger")).
		WithFile(".golangci.yml", source.File(".golangci.yml"))
}

// withSqlc installs sqlc binary in the container.
func (m *Thetree) withSqlc(c *dagger.Container) *dagger.Container {
	return c.WithExec([]string{
		"sh", "-c",
		fmt.Sprintf(
			"curl -sSL https://downloads.sqlc.dev/sqlc_%s_linux_amd64.tar.gz -o /tmp/sqlc.tar.gz && "+
				"tar -xzf /tmp/sqlc.tar.gz -C /usr/local/bin sqlc && chmod +x /usr/local/bin/sqlc",
			sqlcVersion,
		),
	})
}

// withGolangciLint installs golangci-lint binary in the container.
func (m *Thetree) withGolangciLint(c *dagger.Container) *dagger.Container {
	return c.WithExec([]string{
		"sh", "-c",
		fmt.Sprintf(
			"curl -sSL https://github.com/golangci/golangci-lint/releases/download/v%s/golangci-lint-%s-linux-amd64.tar.gz -o /tmp/gcl.tar.gz && "+
				"tar -xzf /tmp/gcl.tar.gz --strip-components=1 -C /usr/local/bin golangci-lint-%s-linux-amd64/golangci-lint && "+
				"chmod +x /usr/local/bin/golangci-lint",
			golangciLintVersion, golangciLintVersion, golangciLintVersion,
		),
	})
}

// withGotestsum installs gotestsum binary in the container.
func (m *Thetree) withGotestsum(c *dagger.Container) *dagger.Container {
	return c.WithExec([]string{
		"sh", "-c",
		fmt.Sprintf(
			"curl -sSL https://github.com/gotestyourself/gotestsum/releases/download/v%s/gotestsum_%s_linux_amd64.tar.gz -o /tmp/gts.tar.gz && "+
				"tar -xzf /tmp/gts.tar.gz -C /usr/local/bin gotestsum && chmod +x /usr/local/bin/gotestsum",
			gotestsumVersion, gotestsumVersion,
		),
	})
}

// sqlcGenerate runs sqlc generate for a service if it has a sqlc.yaml.
func (m *Thetree) sqlcGenerate(c *dagger.Container, service string) *dagger.Container {
	return c.WithExec([]string{
		"sh", "-c",
		fmt.Sprintf(
			"if [ -f api/%s/sqlc.yaml ]; then cd api/%s && sqlc generate; fi",
			service, service,
		),
	})
}

// withOapiCodegen installs the oapi-codegen binary in the container.
// Installing from source (matching `just generate-oapi`'s `go run @v2.8.0`):
// the oapi-codegen/oapi-codegen fork publishes no release binaries, so the
// tarball download fails with a non-gzip 404 page. The production
// Dockerfiles copy the binary from the go-tools image instead.
func (m *Thetree) withOapiCodegen(c *dagger.Container) *dagger.Container {
	return c.WithExec([]string{
		"sh", "-c",
		fmt.Sprintf(
			"GOBIN=/usr/local/bin go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v%s",
			oapiCodegenVersion,
		),
	})
}

// oapiGenerate regenerates the OpenAPI handler bindings (internal/openapi)
// from api-spec.yml for a service. Generated bindings import
// github.com/oapi-codegen/runtime; since the code is not committed (and
// `go mod tidy` would drop an unimported require), the dep is pinned here
// at generation time only — go.mod in the repo stays clean.
func (m *Thetree) oapiGenerate(c *dagger.Container, service string) *dagger.Container {
	return c.WithExec([]string{
		"sh", "-c",
		fmt.Sprintf(
			"cd api/%s && "+
				"oapi-codegen --config oapi-codegen.yaml -generate types -o internal/openapi/types.gen.go api-spec.yml && "+
				"oapi-codegen --config oapi-codegen.yaml -generate chi-server,strict-server -o internal/openapi/server.gen.go api-spec.yml && "+
				"go get github.com/oapi-codegen/runtime@v1.6.0",
			service,
		),
	})
}

// Compile builds a single service.
func (m *Thetree) Compile(ctx context.Context, source *dagger.Directory, service string) (string, error) {
	c := m.withSource(m.withOapiCodegen(m.withSqlc(m.goBase())), source, service)
	c = m.sqlcGenerate(c, service)
	c = m.oapiGenerate(c, service)
	c = c.WithExec([]string{"sh", "-c", fmt.Sprintf("cd api/%s && go build ./...", service)})
	return c.Stdout(ctx)
}

// Lint runs golangci-lint for a single service.
func (m *Thetree) Lint(ctx context.Context, source *dagger.Directory, service string) (string, error) {
	c := m.withSource(m.withGolangciLint(m.withOapiCodegen(m.withSqlc(m.goBase()))), source, service)
	c = m.sqlcGenerate(c, service)
	c = m.oapiGenerate(c, service)
	c = c.WithExec([]string{
		"sh", "-c",
		fmt.Sprintf("cd api/%s && golangci-lint run ./... --config=/workspace/.golangci.yml", service),
	})
	return c.Stdout(ctx)
}

// Test runs gotestsum for a single service.
//
// DB-backed integration tests (lib/testdb, testcontainers) auto-skip here:
// the dagger module runtime has no Docker daemon access, and testdb skips
// when no provider is reachable. CI runs those integration tests on the
// runner, which has Docker — see .forgejo/workflows/ci.yml.
func (m *Thetree) Test(ctx context.Context, source *dagger.Directory, service string) (string, error) {
	c := m.withSource(m.withGotestsum(m.withOapiCodegen(m.withSqlc(m.goBase()))), source, service)
	c = m.sqlcGenerate(c, service)
	c = m.oapiGenerate(c, service)
	c = c.WithExec([]string{
		"sh", "-c",
		fmt.Sprintf("cd api/%s && gotestsum --format testdox --format-hide-empty-pkg ./...", service),
	})
	return c.Stdout(ctx)
}

// CI runs compile, lint, and test for the given services (comma-separated
// or "all"), one service per goroutine — the dagger engine schedules the
// pipelines in parallel, so multi-service runs finish in roughly the time
// of the slowest service instead of the sum.
//
// The test step runs the unit tests; the DB-backed integration tests
// auto-skip (no Docker in the module runtime) and are executed on the CI
// runner by the workflow's `tests` job instead.
// Generated OpenAPI bindings (internal/openapi) are not committed; every
// build path regenerates them from api-spec.yml (see oapiGenerate).
func (m *Thetree) CI(ctx context.Context, source *dagger.Directory, services string) (string, error) {
	list := parseServices(services)
	var (
		mu  sync.Mutex
		out strings.Builder
		wg  sync.WaitGroup
		errs = make(chan error, len(list))
	)
	for _, s := range list {
		s := s
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := m.Compile(ctx, source, s); err != nil {
				errs <- fmt.Errorf("compile %s: %w", s, err)
				return
			}
			if _, err := m.Lint(ctx, source, s); err != nil {
				errs <- fmt.Errorf("lint %s: %w", s, err)
				return
			}
			res, err := m.Test(ctx, source, s)
			if err != nil {
				errs <- fmt.Errorf("test %s: %w", s, err)
				return
			}
			mu.Lock()
			out.WriteString(fmt.Sprintf("--- %s ---\n%s\n", s, res))
			mu.Unlock()
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			return "", err
		}
	}
	return out.String(), nil
}

const registry = "git.trieoh.com"

// Publish builds a service image from its Dockerfile and pushes it to the registry.
// tag is in the format "<service>/v<version>" (e.g. "identityx/v1.2.3").
func (m *Thetree) Publish(
	ctx context.Context,
	source *dagger.Directory,
	registryUsername string,
	registryPassword *dagger.Secret,
	tag string,
) (string, error) {
	parts := strings.SplitN(tag, "/", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid tag %q, expected <service>/v<version>", tag)
	}
	service, version := parts[0], parts[1]

	if strings.HasSuffix(service, "-ui") || strings.Contains(service, "sdk") {
		return "", nil
	}

	img := source.DockerBuild(dagger.DirectoryDockerBuildOpts{
		Dockerfile: fmt.Sprintf("api/%s/Dockerfile", service),
	})

	versionAddr := fmt.Sprintf("%s/trieoh/%s:%s", registry, service, version)
	digest, err := img.
		WithRegistryAuth(registry, registryUsername, registryPassword).
		Publish(ctx, versionAddr)
	if err != nil {
		return "", fmt.Errorf("publish %s: %w", versionAddr, err)
	}

	latestAddr := fmt.Sprintf("%s/trieoh/%s:latest", registry, service)
	if _, err := img.
		WithRegistryAuth(registry, registryUsername, registryPassword).
		Publish(ctx, latestAddr); err != nil {
		return "", fmt.Errorf("publish latest %s: %w", latestAddr, err)
	}

	return digest, nil
}

// FrontendLintTsc runs lint and TypeScript checks for the given frontend
// services (comma-separated or "all").
func (m *Thetree) FrontendLintTsc(
	ctx context.Context,
	// +ignore=["node_modules", "**/node_modules", "dist", "**/dist", ".git", ".wrangler", "**/.wrangler", ".env", "*.env"]
	source *dagger.Directory,
	services string,
) (string, error) {
	list := parseServices(services)
	c := dag.Container().
		From("node:24-bookworm").
		WithMountedCache("/root/.local/share/pnpm/store", dag.CacheVolume("pnpm-store")).
		WithDirectory("/workspace", source).
		WithWorkdir("/workspace").
		WithEnvVariable("CI", "true").
		WithExec([]string{"corepack", "enable"}).
		WithExec([]string{"corepack", "prepare", "pnpm@10", "--activate"}).
		WithExec([]string{"pnpm", "install", "--frozen-lockfile"})

	for _, service := range list {
		c = c.WithExec([]string{"pnpm", "-F", service, "lint"})
		c = c.WithExec([]string{"pnpm", "-F", service, "tsc"})
	}

	return c.Stdout(ctx)
}

// Trivy scans the given project paths (comma-separated, or "all"/empty for
// the whole repo root) with trivy fs: all severities (LOW through
// CRITICAL), secrets, and misconfigurations. It fails on any finding (exit
// code 1), matching the old shell-based trivy-scan.yml. The vulnerability
// DB is cached in a volume, so it is downloaded once per runner and reused
// across runs.
func (m *Thetree) Trivy(ctx context.Context, source *dagger.Directory, projects string) (string, error) {
	c := dag.Container().
		From(fmt.Sprintf("aquasec/trivy:%s", trivyVersion)).
		WithUser("root").
		WithEnvVariable("TRIVY_CACHE_DIR", "/root/.cache/trivy").
		WithMountedCache("/root/.cache/trivy", dag.CacheVolume("trivy-db")).
		WithDirectory("/workspace", source).
		WithWorkdir("/workspace")

	list := []string{"."}
	if projects != "" && projects != "all" {
		list = parseServices(projects) // generic comma-split, see below
	}
	for _, p := range list {
		c = c.WithExec([]string{
			"trivy", "fs",
			"--severity", "LOW,MEDIUM,HIGH,CRITICAL",
			"--scanners", "vuln,secret,misconfig",
			"--exit-code", "1",
			p,
		})
	}
	return c.Stdout(ctx)
}

// parseServices splits a comma-separated list; empty or "all" means every
// service. Despite the name it is a generic comma-splitter (also used for
// trivy project paths).
func parseServices(services string) []string {
	if services == "" || services == "all" {
		return append([]string{}, servicesConst()...)
	}
	// naive split, no external deps
	var out []string
	cur := ""
	for _, r := range services {
		if r == ',' {
			if cur != "" {
				out = append(out, cur)
			}
			cur = ""
			continue
		}
		cur += string(r)
	}
	if cur != "" {
		out = append(out, cur)
	}
	return out
}

func servicesConst() []string {
	return services
}
