package handlers

import (
	"context"

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
	return openapi.ListOccurrenceParticipants200JSONResponse{Items: participants}, nil
}
