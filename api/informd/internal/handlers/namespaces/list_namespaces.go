package namespaces

import (
	"context"
	"time"

	"Informd/internal/openapi"
)

func (h *Handlers) ListNamespaces(ctx context.Context, _ openapi.ListNamespacesRequestObject) (openapi.ListNamespacesResponseObject, error) {
	namespaces, err := h.ops.ListNamespaces(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.ListNamespaces200JSONResponse{
		Code: 200, Data: &namespaces, Timestamp: time.Now(), Module: module,
	}, nil
}
