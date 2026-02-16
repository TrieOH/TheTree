package outbounds

import (
	"GoAuth/internal/domain/scopes"
	"context"

	"github.com/google/uuid"
)

type ScopeRepository interface {
	Create(ctx context.Context, toCreate scopes.Scope) (*scopes.Scope, error)
	GetByIDInternal(ctx context.Context, id uuid.UUID) (*scopes.Scope, error)
	GetByIDExternal(ctx context.Context, id, projectID uuid.UUID) (*scopes.Scope, error)
	GetRootByProjectID(ctx context.Context, projectID uuid.UUID) (*scopes.Scope, error)
	GetProjectScopes(ctx context.Context, projectID uuid.UUID) ([]scopes.Scope, error)

	// Hierarchy methods
	GetAncestors(ctx context.Context, scopeID uuid.UUID) ([]scopes.Scope, error)
	GetGlobalScope(ctx context.Context) (*scopes.Scope, error)
	GetChildren(ctx context.Context, scopeID uuid.UUID) ([]scopes.Scope, error)
}
