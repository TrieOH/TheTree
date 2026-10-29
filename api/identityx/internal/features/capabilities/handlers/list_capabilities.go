package handlers

import (
	"context"

	"IdentityX/internal/openapi"
)

func (h *Handlers) ListCapabilities(ctx context.Context, req openapi.ListCapabilitiesRequestObject) (openapi.ListCapabilitiesResponseObject, error) {
	caps, err := h.ops.List(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.ListCapabilities200JSONResponse{Items: caps}, nil
}
