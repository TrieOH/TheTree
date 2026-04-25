package authz

import (
	"Informd/internal/shared/contracts"
	"context"

	"github.com/MintzyG/FastUtilitiesNet"
)

const ProjectContextKey contextKey = "project"

func WithProject(ctx context.Context, project *contracts.Namespace) context.Context {
	return context.WithValue(ctx, ProjectContextKey, project)
}

func RequireProject(ctx context.Context) (*contracts.Namespace, error) {
	val := ctx.Value(ProjectContextKey)
	if val == nil {
		return nil, fun.ErrInternal("project not found in context")
	}

	ws, ok := val.(*contracts.Namespace)
	if !ok {
		return nil, fun.Errf("Invalid project type, was: %T", val).Internal()
	}

	return ws, nil
}
