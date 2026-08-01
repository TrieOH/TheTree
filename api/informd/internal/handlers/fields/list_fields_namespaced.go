package fields

import (
	"context"
	"time"

	"Informd/internal/openapi"
)

func (h *Handlers) ListFieldsNamespaced(ctx context.Context, req openapi.ListFieldsNamespacedRequestObject) (openapi.ListFieldsNamespacedResponseObject, error) {
	fields, err := h.ops.ListNamespaced(ctx, req.FormId, req.NamespaceId, req.StepId)
	if err != nil {
		return nil, err
	}
	return openapi.ListFieldsNamespaced200JSONResponse{
		Code: 200, Data: &fields, Timestamp: time.Now(), Module: module,
	}, nil
}
