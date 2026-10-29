package handlers

import (
	"context"

	"Informd/internal/openapi"
)

func (h *Handlers) ArchiveFormNamespaced(ctx context.Context, req openapi.ArchiveFormNamespacedRequestObject) (openapi.ArchiveFormNamespacedResponseObject, error) {
	form, err := h.ops.ArchiveForm(ctx, req.NamespaceId, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.ArchiveFormNamespaced200JSONResponse(*form), nil
}
