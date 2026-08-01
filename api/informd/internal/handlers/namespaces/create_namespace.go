package namespaces

import (
	"context"
	"time"

	"Informd/internal/openapi"
)

func (h *Handlers) CreateNamespace(ctx context.Context, req openapi.CreateNamespaceRequestObject) (openapi.CreateNamespaceResponseObject, error) {
	namespace, err := h.ops.Create(ctx, req.Body.Name)
	if err != nil {
		return nil, err
	}
	return openapi.CreateNamespace201JSONResponse{
		Code: 201, Data: namespace, Timestamp: time.Now(), Module: module,
	}, nil
}
