package handlers

import (
	"context"

	"Informd/internal/openapi"
)

func (h *Handlers) CreateNamespaceForm(ctx context.Context, req openapi.CreateNamespaceFormRequestObject) (openapi.CreateNamespaceFormResponseObject, error) {
	form, err := h.ops.CreateForm(ctx, req.Body.Title, req.NamespaceId)
	if err != nil {
		return nil, err
	}
	return openapi.CreateNamespaceForm201JSONResponse(*form), nil
}
