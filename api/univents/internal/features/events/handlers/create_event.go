package handlers

import (
	"context"

	"univents/internal/openapi"
	"univents/models"
)

func (h *Handlers) CreateEvent(ctx context.Context, req openapi.CreateEventRequestObject) (openapi.CreateEventResponseObject, error) {
	event, err := h.ops.Create(ctx, models.CreateEventInput{
		FullName:     req.Body.FullName,
		Acronym:      req.Body.Acronym,
		Slug:         req.Body.Slug,
		Description:  req.Body.Description,
		ContactEmail: req.Body.ContactEmail,
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateEvent201JSONResponse(*event), nil
}
