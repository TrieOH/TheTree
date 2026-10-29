package handlers

import (
	"context"

	"Informd/internal/openapi"
)

func (h *Handlers) ListStepsNamespaced(ctx context.Context, req openapi.ListStepsNamespacedRequestObject) (openapi.ListStepsNamespacedResponseObject, error) {
	steps, err := h.ops.ListNamespaced(ctx, req.FormId, req.NamespaceId)
	if err != nil {
		return nil, err
	}
	return openapi.ListStepsNamespaced200JSONResponse{Items: steps}, nil
}
