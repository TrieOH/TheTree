package namespaces

import (
	"context"
	"time"

	"Informd/internal/openapi"
)

func (h *Handlers) GetFullFormNamespaced(ctx context.Context, req openapi.GetFullFormNamespacedRequestObject) (openapi.GetFullFormNamespacedResponseObject, error) {
	form, err := h.ops.GetFullForm(ctx, req.NamespaceId, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.GetFullFormNamespaced200JSONResponse{
		Code: 200, Data: form, Timestamp: time.Now(), Module: module,
	}, nil
}
