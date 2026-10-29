package handlers

import (
	"context"

	"Informd/internal/openapi"
)

func (h *Handlers) GetFormResponseCount(ctx context.Context, req openapi.GetFormResponseCountRequestObject) (openapi.GetFormResponseCountResponseObject, error) {
	count, err := h.ops.GetResponseCount(ctx, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.GetFormResponseCount200JSONResponse{Count: count}, nil
}
