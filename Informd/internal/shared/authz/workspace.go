package authz

import (
	"TrieForms/internal/shared/errx"
	"TrieForms/internal/shared/types"
	"context"
	"fmt"
)

const ProjectContextKey contextKey = "project"

func WithProject(ctx context.Context, project *types.Project) context.Context {
	return context.WithValue(ctx, ProjectContextKey, project)
}

func RequireProject(ctx context.Context) (*types.Project, error) {
	val := ctx.Value(ProjectContextKey)
	if val == nil {
		return nil, errx.NotFound("project").SetMessage("project not found in context")
	}

	ws, ok := val.(*types.Project)
	if !ok {
		return nil, errx.Invalid("project").SetMessage(fmt.Sprintf("type was %T", val))
	}

	return ws, nil
}
