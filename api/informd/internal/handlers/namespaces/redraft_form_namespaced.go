package namespaces

import (
	"context"
	"time"

	"Informd/internal/openapi"
)

func (h *Handlers) RedraftFormNamespaced(ctx context.Context, req openapi.RedraftFormNamespacedRequestObject) (openapi.RedraftFormNamespacedResponseObject, error) {
	form, err := h.ops.ReDraftForm(ctx, req.NamespaceId, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.RedraftFormNamespaced200JSONResponse{
		Code: 200, Data: form, Timestamp: time.Now(), Module: module,
	}, nil
}
