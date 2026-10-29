package handlers

import (
	"context"

	"Informd/internal/openapi"
)

func (h *Handlers) GetFullForm(ctx context.Context, req openapi.GetFullFormRequestObject) (openapi.GetFullFormResponseObject, error) {
	form, err := h.ops.GetFull(ctx, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.GetFullForm200JSONResponse(*form), nil
}
