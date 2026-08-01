package namespaces

import (
	"context"
	"time"

	"Informd/internal/openapi"
)

func (h *Handlers) ListNamespaceForms(ctx context.Context, req openapi.ListNamespaceFormsRequestObject) (openapi.ListNamespaceFormsResponseObject, error) {
	forms, err := h.ops.ListForms(ctx, req.NamespaceId)
	if err != nil {
		return nil, err
	}
	return openapi.ListNamespaceForms200JSONResponse{
		Code: 200, Data: &forms, Timestamp: time.Now(), Module: module,
	}, nil
}
