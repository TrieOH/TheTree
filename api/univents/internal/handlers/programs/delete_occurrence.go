package programs

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) DeleteOccurrence(ctx context.Context, req openapi.DeleteOccurrenceRequestObject) (openapi.DeleteOccurrenceResponseObject, error) {
	occurrence, err := h.ops.DeleteOccurrence(ctx, req.OccurrenceId)
	if err != nil {
		return nil, err
	}
	return openapi.DeleteOccurrence200JSONResponse{
		Code: 200, Data: occurrence, Timestamp: time.Now(), Module: module,
	}, nil
}
