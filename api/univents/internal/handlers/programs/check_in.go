package programs

import (
	"context"
	"time"

	idx "sdk/identityx"

	"univents/internal/openapi"
)

func (h *Handlers) CheckInOccurrence(ctx context.Context, req openapi.CheckInOccurrenceRequestObject) (openapi.CheckInOccurrenceResponseObject, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	part, err := h.ops.CheckIn(ctx, req.OccurrenceId, req.Body.AttendeeId, ident.Sub.ID)
	if err != nil {
		return nil, err
	}
	return openapi.CheckInOccurrence200JSONResponse{
		Code: 200, Data: part, Timestamp: time.Now(), Module: module,
	}, nil
}
