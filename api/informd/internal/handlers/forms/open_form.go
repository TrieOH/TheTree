package forms

import (
	"context"
	"time"

	"Informd/internal/openapi"
)

func (h *Handlers) OpenForm(ctx context.Context, req openapi.OpenFormRequestObject) (openapi.OpenFormResponseObject, error) {
	form, err := h.ops.Open(ctx, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.OpenForm200JSONResponse{
		Code: 200, Data: form, Timestamp: time.Now(), Module: module,
	}, nil
}
