// A generated module for Thetree functions

package main

import (
	"context"
	"dagger/thetree/internal/dagger"
	"fmt"
	"strings"
)

type Thetree struct{}

var services = []string{"identityx", "informd", "payssage", "univents"}

const (
	goVersion           = "1.26.4"
	sqlcVersion         = "1.31.1"
	gotestsumVersion    = "1.13.0"
	golangciLintVersion = "2.12.2"
)

// baseGo returns a container with the full repo source mounted.
func (m *Thetree) baseGo(source *dagger.Directory) *dagger.Container {
	return m.baseGoScoped(source, "")
}

// baseGoScoped returns a container with only the workspace modules needed
// by the given service. Pass an empty service string for the full source.
func (m *Thetree) baseGoScoped(source *dagger.Directory, service string) *dagger.Container {
	dirs := m.scopedSource(source, service)
	return dag.Container().
		From(fmt.Sprintf("golang:%s-bookworm", goVersion)).
		WithMountedCache("/root/go/pkg/mod", dag.CacheVolume("go-mod")).
		WithMountedCache("/root/.cache/go-build", dag.CacheVolume("go-build")).
		WithDirectory("/workspace", dirs).
		WithWorkdir("/workspace").
		WithEnvVariable("GOWORK", "/workspace/go.work")
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

// Compile builds a single service.
func (m *Thetree) Compile(ctx context.Context, source *dagger.Directory, service string) (string, error) {
	c := m.withSqlc(m.baseGoScoped(source, service))
	c = m.sqlcGenerate(c, service)
	c = c.WithExec([]string{"sh", "-c", fmt.Sprintf("cd api/%s && go build ./...", service)})
	return c.Stdout(ctx)
}

// Lint runs golangci-lint for a single service.
func (m *Thetree) Lint(ctx context.Context, source *dagger.Directory, service string) (string, error) {
	c := m.withGolangciLint(m.withSqlc(m.baseGoScoped(source, service)))
	c = m.sqlcGenerate(c, service)
	c = c.WithExec([]string{
		"sh", "-c",
		fmt.Sprintf("cd api/%s && golangci-lint run ./... --config=/workspace/.golangci.yml", service),
	})
	return c.Stdout(ctx)
}

// Test runs gotestsum for a single service.
func (m *Thetree) Test(ctx context.Context, source *dagger.Directory, service string) (string, error) {
	c := m.withGotestsum(m.withSqlc(m.baseGoScoped(source, service)))
	c = m.sqlcGenerate(c, service)
	c = c.WithExec([]string{
		"sh", "-c",
		fmt.Sprintf("cd api/%s && gotestsum --format testdox --format-hide-empty-pkg ./...", service),
	})
	return c.Stdout(ctx)
}

// CI runs compile, lint, and test for the given services (comma-separated or "all").
func (m *Thetree) CI(ctx context.Context, source *dagger.Directory, services string) (string, error) {
	list := parseServices(services)
	var out string
	for _, s := range list {
		if _, err := m.Compile(ctx, source, s); err != nil {
			return "", fmt.Errorf("compile %s: %w", s, err)
		}
		if _, err := m.Lint(ctx, source, s); err != nil {
			return "", fmt.Errorf("lint %s: %w", s, err)
		}
		res, err := m.Test(ctx, source, s)
		if err != nil {
			return "", fmt.Errorf("test %s: %w", s, err)
		}
		out += fmt.Sprintf("--- %s ---\n%s\n", s, res)
	}
	return out, nil
}

const registry = "git.trieoh.com"

// Publish builds a service image from its Dockerfile and pushes it to the registry.
// tag is in the format "<service>/v<version>" (e.g. "identityx/v1.2.3").
func (m *Thetree) Publish(
	ctx context.Context,
	source *dagger.Directory,
	registryUsername string,
	registryPassword string,
	tag string,
) (string, error) {
	parts := strings.SplitN(tag, "/", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid tag %q, expected <service>/v<version>", tag)
	}
	service, version := parts[0], parts[1]

	img := source.DockerBuild(dagger.DirectoryDockerBuildOpts{
		Dockerfile: fmt.Sprintf("api/%s/Dockerfile", service),
	})

	secret := dag.SetSecret("registry-password", registryPassword)

	versionAddr := fmt.Sprintf("%s/trieoh/%s:%s", registry, service, version)
	digest, err := img.
		WithRegistryAuth(registry, registryUsername, secret).
		Publish(ctx, versionAddr)
	if err != nil {
		return "", fmt.Errorf("publish %s: %w", versionAddr, err)
	}

	latestAddr := fmt.Sprintf("%s/trieoh/%s:latest", registry, service)
	if _, err := img.
		WithRegistryAuth(registry, registryUsername, secret).
		Publish(ctx, latestAddr); err != nil {
		return "", fmt.Errorf("publish latest %s: %w", latestAddr, err)
	}

	return digest, nil
}

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
