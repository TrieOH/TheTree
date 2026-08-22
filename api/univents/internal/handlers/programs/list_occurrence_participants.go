package programs

import (
	"context"
	"time"

	idx "sdk/identityx"

	"univents/internal/openapi"
)

func (h *Handlers) ListOccurrenceParticipants(ctx context.Context, req openapi.ListOccurrenceParticipantsRequestObject) (openapi.ListOccurrenceParticipantsResponseObject, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	participants, err := h.ops.Participants(ctx, req.OccurrenceId, ident.Sub.ID)
	if err != nil {
		return nil, err
	}
	return openapi.ListOccurrenceParticipants200JSONResponse{
		Code: 200, Data: &participants, Timestamp: time.Now(), Module: module,
	}, nil
}
