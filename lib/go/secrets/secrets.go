// Package secrets is the seam between a Backend's configuration and wherever
// its secret values live. A Backend reads its config from the environment
// (caarlos0/env); a Source only fills the environment in before that parse,
// so no config struct, and no code past boot, knows which secret manager —
// if any — is behind it. Swapping managers is an adapter, never a refactor.
//
// Adapters: Env (the process environment already holds everything, e.g.
// compose env_file) and awssm (AWS Secrets Manager). provider.FromEnv picks
// one from SECRETS_PROVIDER.
package secrets

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Source loads a flat key → value set of secrets.
type Source interface {
	Load(ctx context.Context) (map[string]string, error)
}

// Store is a Source that can also be written — what `just inject` seeds.
type Store interface {
	Source
	Put(ctx context.Context, values map[string]string) error
}

// Apply loads src and exports every value whose key is not already set in
// the process environment. Explicit environment always wins, so a function
// or container can override a single value without touching the secret.
// It returns the keys it exported.
func Apply(ctx context.Context, src Source) ([]string, error) {
	values, err := src.Load(ctx)
	if err != nil {
		return nil, fmt.Errorf("load secrets: %w", err)
	}
	applied := make([]string, 0, len(values))
	for key, value := range values {
		if _, set := os.LookupEnv(key); set {
			continue
		}
		err = os.Setenv(key, value)
		if err != nil {
			return nil, fmt.Errorf("export secret %s: %w", key, err)
		}
		applied = append(applied, key)
	}
	sort.Strings(applied)
	return applied, nil
}

// Env is the no-op Source: the process environment is the source of truth
// (compose env_file, a systemd EnvironmentFile, a vault agent template).
type Env struct{}

func (Env) Load(context.Context) (map[string]string, error) { return map[string]string{}, nil }

// ParseDotenv parses KEY=VALUE lines: blank lines and # comments are
// skipped, an optional `export ` prefix is dropped, and one level of
// matching single or double quotes is stripped from the value. Later files
// win over earlier ones when merged by the caller.
func ParseDotenv(path string) (map[string]string, error) {
	f, err := os.Open(path) //nolint:gosec // operator-supplied path
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	values := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		value = strings.TrimSpace(value)
		if len(value) >= 2 && value[0] == value[len(value)-1] && (value[0] == '"' || value[0] == '\'') {
			value = value[1 : len(value)-1]
		}
		values[strings.TrimSpace(key)] = value
	}
	return values, scanner.Err()
}
