package authz

import (
	"context"
	"fmt"
	"payssage/internal/shared/contracts"
	"payssage/internal/shared/errx"
)

const WorkspaceContextKey contextKey = "workspace"

func WithWorkspace(ctx context.Context, workspace *contracts.Workspace) context.Context {
	return context.WithValue(ctx, WorkspaceContextKey, workspace)
}

func RequireWorkspace(ctx context.Context) (*contracts.Workspace, error) {
	val := ctx.Value(WorkspaceContextKey)
	if val == nil {
		return nil, errx.NotFound("workspace").SetMessage("workspace not found in context")
	}

	ws, ok := val.(*contracts.Workspace)
	if !ok {
		return nil, errx.Invalid("workspace").SetMessage(fmt.Sprintf("type was %T", val))
	}

	return ws, nil
}
