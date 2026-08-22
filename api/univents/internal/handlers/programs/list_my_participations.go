package programs

import (
	"context"
	"time"

	idx "sdk/identityx"

	"univents/internal/openapi"
)

func (h *Handlers) ListMyParticipations(ctx context.Context, req openapi.ListMyParticipationsRequestObject) (openapi.ListMyParticipationsResponseObject, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	participations, err := h.ops.MyParticipations(ctx, req.EditionId, ident.Sub.ID)
	if err != nil {
		return nil, err
	}
	return openapi.ListMyParticipations200JSONResponse{
		Code: 200, Data: &participations, Timestamp: time.Now(), Module: module,
	}, nil
}
