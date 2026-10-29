package handlers

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) GetEventBySlug(ctx context.Context, req openapi.GetEventBySlugRequestObject) (openapi.GetEventBySlugResponseObject, error) {
	event, err := h.ops.GetBySlug(ctx, req.EventSlug)
	if err != nil {
		return nil, err
	}
	return openapi.GetEventBySlug200JSONResponse(*event), nil
}
