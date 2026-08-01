package forms

import (
	"context"
	"time"

	"Informd/internal/openapi"
)

func (h *Handlers) GetFullForm(ctx context.Context, req openapi.GetFullFormRequestObject) (openapi.GetFullFormResponseObject, error) {
	form, err := h.ops.GetFull(ctx, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.GetFullForm200JSONResponse{
		Code: 200, Data: form, Timestamp: time.Now(), Module: module,
	}, nil
}
