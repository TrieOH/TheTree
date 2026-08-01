package forms

import (
	"context"
	"time"

	"Informd/internal/openapi"
)

func (h *Handlers) GetAnswerableForm(ctx context.Context, req openapi.GetAnswerableFormRequestObject) (openapi.GetAnswerableFormResponseObject, error) {
	form, err := h.ops.GetAnswerable(ctx, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.GetAnswerableForm200JSONResponse{
		Code: 200, Data: form, Timestamp: time.Now(), Module: module,
	}, nil
}
