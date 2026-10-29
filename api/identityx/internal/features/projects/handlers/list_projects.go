package handlers

import (
	"context"

	"IdentityX/internal/openapi"
)

func (h *Handlers) ListProjects(ctx context.Context, _ openapi.ListProjectsRequestObject) (openapi.ListProjectsResponseObject, error) {
	projects, err := h.ops.ListProjects(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.ListProjects200JSONResponse{Items: projects}, nil
}
