package programs

import (
	"context"
	"time"

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
	return openapi.CreateProgramOccurrence201JSONResponse{
		Code: 201, Data: occurrence, Timestamp: time.Now(), Module: module,
	}, nil
}
