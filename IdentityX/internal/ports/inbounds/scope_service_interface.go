package inbounds

import (
	"context"
)

type ScopeService interface {
	Create(ctx context.Context, in CreateScopeInput) (*ScopeOutput, error)
	GetByIDExternal(ctx context.Context, in GetScopeInput) (*ScopeOutput, error)
	GetProjectScopesExternal(ctx context.Context, in GetScopeInput) ([]ScopeOutput, error)
}
