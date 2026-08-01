package namespaces

import (
	"context"
	"time"

	"Informd/internal/openapi"
)

func (h *Handlers) GetFormResponseCountNamespaced(ctx context.Context, req openapi.GetFormResponseCountNamespacedRequestObject) (openapi.GetFormResponseCountNamespacedResponseObject, error) {
	count, err := h.ops.GetFormResponseCount(ctx, req.NamespaceId, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.GetFormResponseCountNamespaced200JSONResponse{
		Code: 200,
		Data: &struct {
			Count int `json:"count"`
		}{Count: count},
		Timestamp: time.Now(), Module: module,
	}, nil
}
