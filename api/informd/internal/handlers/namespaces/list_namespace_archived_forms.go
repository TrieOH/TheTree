package namespaces

import (
	"context"
	"time"

	"Informd/internal/openapi"
)

func (h *Handlers) ListNamespaceArchivedForms(ctx context.Context, req openapi.ListNamespaceArchivedFormsRequestObject) (openapi.ListNamespaceArchivedFormsResponseObject, error) {
	forms, err := h.ops.ListArchivedForms(ctx, req.NamespaceId)
	if err != nil {
		return nil, err
	}
	return openapi.ListNamespaceArchivedForms200JSONResponse{
		Code: 200, Data: &forms, Timestamp: time.Now(), Module: module,
	}, nil
}
