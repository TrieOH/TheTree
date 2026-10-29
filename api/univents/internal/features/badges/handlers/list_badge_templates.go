package handlers

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) ListBadgeTemplates(ctx context.Context, req openapi.ListBadgeTemplatesRequestObject) (openapi.ListBadgeTemplatesResponseObject, error) {
	templates, err := h.ops.ListTemplates(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.ListBadgeTemplates200JSONResponse{Items: templates}, nil
}
