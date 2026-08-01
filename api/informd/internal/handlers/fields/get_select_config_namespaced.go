package fields

import (
	"context"
	"time"

	"Informd/internal/openapi"
)

func (h *Handlers) GetSelectConfigNamespaced(ctx context.Context, req openapi.GetSelectConfigNamespacedRequestObject) (openapi.GetSelectConfigNamespacedResponseObject, error) {
	config, err := h.ops.GetSelectConfigNamespaced(ctx, req.FormId, req.NamespaceId, req.FieldId)
	if err != nil {
		return nil, err
	}
	return openapi.GetSelectConfigNamespaced200JSONResponse{
		Code: 200, Data: config, Timestamp: time.Now(), Module: module,
	}, nil
}
