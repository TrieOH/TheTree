package forms

import (
	"context"
	"time"

	"Informd/internal/openapi"
)

func (h *Handlers) ListMyForms(ctx context.Context, _ openapi.ListMyFormsRequestObject) (openapi.ListMyFormsResponseObject, error) {
	forms, err := h.ops.ListForms(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.ListMyForms200JSONResponse{
		Code: 200, Data: &forms, Timestamp: time.Now(), Module: module,
	}, nil
}
