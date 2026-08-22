package programs

import (
	"context"
	"time"

	idx "sdk/identityx"

	"univents/internal/openapi"
)

func (h *Handlers) MarkParticipationAttended(ctx context.Context, req openapi.MarkParticipationAttendedRequestObject) (openapi.MarkParticipationAttendedResponseObject, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	part, err := h.ops.MarkAttended(ctx, req.ParticipationId, ident.Sub.ID)
	if err != nil {
		return nil, err
	}
	return openapi.MarkParticipationAttended200JSONResponse{
		Code: 200, Data: part, Timestamp: time.Now(), Module: module,
	}, nil
}
