package handlers

import (
	"context"

	"Informd/internal/openapi"
)

func (h *Handlers) CreateForm(ctx context.Context, req openapi.CreateFormRequestObject) (openapi.CreateFormResponseObject, error) {
	form, err := h.ops.Create(ctx, req.Body.Title)
	if err != nil {
		return nil, err
	}
	return openapi.CreateForm201JSONResponse(*form), nil
}
