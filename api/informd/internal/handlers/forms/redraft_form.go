package forms

import (
	"context"
	"time"

	"Informd/internal/openapi"
)

func (h *Handlers) RedraftForm(ctx context.Context, req openapi.RedraftFormRequestObject) (openapi.RedraftFormResponseObject, error) {
	form, err := h.ops.ReDraft(ctx, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.RedraftForm200JSONResponse{
		Code: 200, Data: form, Timestamp: time.Now(), Module: module,
	}, nil
}
