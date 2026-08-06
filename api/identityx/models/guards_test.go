package models

import (
	"context"
	"testing"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func projectScopedCtx() context.Context {
	pid := uuid.New()
	return WithIdentity(context.Background(), &Identity{
		Sub: Subject{ID: uuid.New(), ProjectID: &pid},
	})
}

func platformCtx() context.Context {
	return WithIdentity(context.Background(), &Identity{
		Sub: Subject{ID: uuid.New(), ProjectID: nil},
	})
}

// TestRequireClientOnly asserts the guard's unit behavior: it passes a
// platform-level identity, rejects a project-scoped one, and rejects an
// unauthenticated context.
func TestRequireClientOnly(t *testing.T) {
	err := RequireClientOnly(platformCtx())
	if err != nil {
		t.Fatalf("platform-level identity must pass the guard, got %v", err)
	}
	err = RequireClientOnly(projectScopedCtx())
	if !fun.Is(err, fun.CodeUnauthorized) {
		t.Fatalf("project-scoped identity must be rejected with 401, got %v", err)
	}
	err = RequireClientOnly(context.Background())
	if !fun.Is(err, fun.CodeUnauthorized) {
		t.Fatalf("unauthenticated context must be rejected with 401, got %v", err)
	}
}

// TestRequireProjectClientOnly asserts the mirror guard (defined for future
// project-scoped operations): it passes a project-scoped identity and
// rejects a platform-level one.
func TestRequireProjectClientOnly(t *testing.T) {
	err := RequireProjectClientOnly(projectScopedCtx())
	if err != nil {
		t.Fatalf("project-scoped identity must pass the guard, got %v", err)
	}
	err = RequireProjectClientOnly(platformCtx())
	if !fun.Is(err, fun.CodeUnauthorized) {
		t.Fatalf("platform-level identity must be rejected with 401, got %v", err)
	}
	err = RequireProjectClientOnly(context.Background())
	if !fun.Is(err, fun.CodeUnauthorized) {
		t.Fatalf("unauthenticated context must be rejected with 401, got %v", err)
	}
}
