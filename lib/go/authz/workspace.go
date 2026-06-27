package authz

import (
	"context"
	"payssage/models"

	"github.com/MintzyG/fun"
)

type contextKey string

const WorkspaceContextKey contextKey = "workspace"

func WithWorkspace(ctx context.Context, workspace *models.Workspace) context.Context {
	return context.WithValue(ctx, WorkspaceContextKey, workspace)
}

func RequireWorkspace(ctx context.Context) (*models.Workspace, error) {
	val := ctx.Value(WorkspaceContextKey)
	if val == nil {
		return nil, fun.ErrNotFound("workspace not found in context")
	}

	ws, ok := val.(*models.Workspace)
	if !ok {
		return nil, fun.Errf("invalid workspace type, was: %T", val).Internal()
	}

	return ws, nil
}
