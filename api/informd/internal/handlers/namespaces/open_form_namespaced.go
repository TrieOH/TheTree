package namespaces

import (
	"context"
	"time"

	"Informd/internal/openapi"
)

func (h *Handlers) OpenFormNamespaced(ctx context.Context, req openapi.OpenFormNamespacedRequestObject) (openapi.OpenFormNamespacedResponseObject, error) {
	form, err := h.ops.OpenForm(ctx, req.NamespaceId, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.OpenFormNamespaced200JSONResponse{
		Code: 200, Data: form, Timestamp: time.Now(), Module: module,
	}, nil
}
