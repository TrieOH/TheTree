package namespaces

import (
	"context"
	"time"

	"Informd/internal/openapi"
)

func (h *Handlers) CloseFormNamespaced(ctx context.Context, req openapi.CloseFormNamespacedRequestObject) (openapi.CloseFormNamespacedResponseObject, error) {
	form, err := h.ops.CloseForm(ctx, req.NamespaceId, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.CloseFormNamespaced200JSONResponse{
		Code: 200, Data: form, Timestamp: time.Now(), Module: module,
	}, nil
}
