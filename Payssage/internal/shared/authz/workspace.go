package authz

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/errx"
	"context"
	"fmt"
)

const WorkspaceContextKey contextKey = "workspace"

func WithWorkspace(ctx context.Context, workspace *domain.Workspace) context.Context {
	return context.WithValue(ctx, WorkspaceContextKey, workspace)
}

func RequireWorkspace(ctx context.Context) (*domain.Workspace, error) {
	val := ctx.Value(WorkspaceContextKey)
	if val == nil {
		return nil, errx.NotFound("workspace").SetMessage("workspace not found in context")
	}

	ws, ok := val.(*domain.Workspace)
	if !ok {
		return nil, errx.Invalid("workspace").SetMessage(fmt.Sprintf("type was %T", val))
	}

	return ws, nil
}
