package programs

import (
	"context"
	"time"

	"univents/internal/openapi"
	"univents/models"
)

func (h *Handlers) PatchOccurrence(ctx context.Context, req openapi.PatchOccurrenceRequestObject) (openapi.PatchOccurrenceResponseObject, error) {
	occurrence, err := h.ops.PatchOccurrence(ctx, models.PatchProgramOccurrenceInput{
		OccurrenceID: req.OccurrenceId,
		StartsAt:     req.Body.StartsAt,
		EndsAt:       req.Body.EndsAt,
		MaxCapacity:  req.Body.MaxCapacity,
	})
	if err != nil {
		return nil, err
	}
	return openapi.PatchOccurrence200JSONResponse{
		Code: 200, Data: occurrence, Timestamp: time.Now(), Module: module,
	}, nil
}
