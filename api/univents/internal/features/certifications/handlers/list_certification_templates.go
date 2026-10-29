package handlers

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) ListCertificationTemplates(ctx context.Context, req openapi.ListCertificationTemplatesRequestObject) (openapi.ListCertificationTemplatesResponseObject, error) {
	templates, err := h.ops.ListTemplates(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListCertificationTemplates200JSONResponse{Items: templates}, nil
}
