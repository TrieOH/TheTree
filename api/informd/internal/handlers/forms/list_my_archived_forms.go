package forms

import (
	"context"
	"time"

	"Informd/internal/openapi"
)

func (h *Handlers) ListMyArchivedForms(ctx context.Context, _ openapi.ListMyArchivedFormsRequestObject) (openapi.ListMyArchivedFormsResponseObject, error) {
	forms, err := h.ops.ListArchivedForms(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.ListMyArchivedForms200JSONResponse{
		Code: 200, Data: &forms, Timestamp: time.Now(), Module: module,
	}, nil
}
