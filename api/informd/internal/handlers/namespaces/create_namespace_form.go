package namespaces

import (
	"context"
	"time"

	"Informd/internal/openapi"
)

func (h *Handlers) CreateNamespaceForm(ctx context.Context, req openapi.CreateNamespaceFormRequestObject) (openapi.CreateNamespaceFormResponseObject, error) {
	form, err := h.ops.CreateForm(ctx, req.Body.Title, req.NamespaceId)
	if err != nil {
		return nil, err
	}
	return openapi.CreateNamespaceForm201JSONResponse{
		Code: 201, Data: form, Timestamp: time.Now(), Module: module,
	}, nil
}
