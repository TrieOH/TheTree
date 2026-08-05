package forms

import (
	"context"
	"time"

	"Informd/internal/openapi"
)

func (h *Handlers) CloseForm(ctx context.Context, req openapi.CloseFormRequestObject) (openapi.CloseFormResponseObject, error) {
	form, err := h.ops.Close(ctx, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.CloseForm200JSONResponse{
		Code: 200, Data: form, Timestamp: time.Now(), Module: module,
	}, nil
}
