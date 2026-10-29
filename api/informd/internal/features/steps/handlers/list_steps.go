package handlers

import (
	"context"

	"Informd/internal/openapi"
)

func (h *Handlers) ListSteps(ctx context.Context, req openapi.ListStepsRequestObject) (openapi.ListStepsResponseObject, error) {
	steps, err := h.ops.List(ctx, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.ListSteps200JSONResponse{Items: steps}, nil
}
