package fields

import (
	"context"
	"time"

	"Informd/internal/openapi"
)

func (h *Handlers) DeleteFieldNamespaced(ctx context.Context, req openapi.DeleteFieldNamespacedRequestObject) (openapi.DeleteFieldNamespacedResponseObject, error) {
	err := h.ops.DeleteNamespaced(ctx, req.NamespaceId, req.FormId, req.FieldId)
	if err != nil {
		return nil, err
	}
	return openapi.DeleteFieldNamespaced200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}
