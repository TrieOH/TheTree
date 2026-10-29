package handlers

import (
	"context"

	"univents/internal/openapi"
	"univents/models"
)

func (h *Handlers) CreateProgramOccurrence(ctx context.Context, req openapi.CreateProgramOccurrenceRequestObject) (openapi.CreateProgramOccurrenceResponseObject, error) {
	occurrence, err := h.ops.CreateOccurrence(ctx, models.CreateProgramOccurrenceInput{
		ProgramID:   req.ProgramId,
		StartsAt:    req.Body.StartsAt,
		EndsAt:      req.Body.EndsAt,
		MaxCapacity: req.Body.MaxCapacity,
	})
	if err != nil {
		return nil, err
	}
	return openapi.CreateProgramOccurrence201JSONResponse(*occurrence), nil
}
