package handlers

import (
	"context"

	"Informd/internal/openapi"
)

func (h *Handlers) CreateNamespace(ctx context.Context, req openapi.CreateNamespaceRequestObject) (openapi.CreateNamespaceResponseObject, error) {
	namespace, err := h.ops.Create(ctx, req.Body.Name)
	if err != nil {
		return nil, err
	}
	return openapi.CreateNamespace201JSONResponse(*namespace), nil
}
