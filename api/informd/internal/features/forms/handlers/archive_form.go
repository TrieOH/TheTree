package handlers

import (
	"context"

	"Informd/internal/openapi"
)

func (h *Handlers) ArchiveForm(ctx context.Context, req openapi.ArchiveFormRequestObject) (openapi.ArchiveFormResponseObject, error) {
	form, err := h.ops.Archive(ctx, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.ArchiveForm200JSONResponse(*form), nil
}
