package authz

import (
	"TrieForms/internal/shared/contracts"
	"context"

	fun "github.com/MintzyG/FastUtilitiesNet/response"
)

const ProjectContextKey contextKey = "project"

func WithProject(ctx context.Context, project *contracts.Project) context.Context {
	return context.WithValue(ctx, ProjectContextKey, project)
}

func RequireProject(ctx context.Context) (*contracts.Project, error) {
	val := ctx.Value(ProjectContextKey)
	if val == nil {
		return nil, fun.NewError("project not found in context").Internal()
	}

	ws, ok := val.(*contracts.Project)
	if !ok {
		return nil, fun.NewErrorf("Invalid project type, was: %T", val).Internal()
	}

	return ws, nil
}
