package handlers

import (
	"IdentityX/internal/openapi"
	"context"
)

func (h *Handlers) ListProjectOAuthProviders(ctx context.Context, req openapi.ListProjectOAuthProvidersRequestObject) (openapi.ListProjectOAuthProvidersResponseObject, error) {
	providers, err := h.ops.ListByProject(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.ListProjectOAuthProviders200JSONResponse{Items: providers}, nil
}
