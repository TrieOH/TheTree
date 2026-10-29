package handlers

import (
	"context"

	"Informd/internal/openapi"
)

func (h *Handlers) GetFormResponseCountNamespaced(ctx context.Context, req openapi.GetFormResponseCountNamespacedRequestObject) (openapi.GetFormResponseCountNamespacedResponseObject, error) {
	count, err := h.ops.GetFormResponseCount(ctx, req.NamespaceId, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.GetFormResponseCountNamespaced200JSONResponse{Count: count}, nil
}
