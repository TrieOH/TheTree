package inbounds

import (
	"context"
)

type ScopeService interface {
	Create(ctx context.Context, in CreateScopeInput) (*ScopeOutput, error)
	UpdateMeta(ctx context.Context, in UpdateProjectScopeMetaInput) error
	GetByIDExternal(ctx context.Context, in GetScopeInput) (*ScopeOutput, error)
	GetProjectScopesExternal(ctx context.Context, in GetScopeInput) ([]ScopeOutput, error)
	Delete(ctx context.Context, in GetScopeInput) error
}
