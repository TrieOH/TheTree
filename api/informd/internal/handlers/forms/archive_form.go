package forms

import (
	"context"
	"time"

	"Informd/internal/openapi"
)

func (h *Handlers) ArchiveForm(ctx context.Context, req openapi.ArchiveFormRequestObject) (openapi.ArchiveFormResponseObject, error) {
	form, err := h.ops.Archive(ctx, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.ArchiveForm200JSONResponse{
		Code: 200, Data: form, Timestamp: time.Now(), Module: module,
	}, nil
}
