package programs

import (
	"context"
	"time"

	idx "sdk/identityx"

	"univents/internal/openapi"
)

func (h *Handlers) DeregisterOccurrence(ctx context.Context, req openapi.DeregisterOccurrenceRequestObject) (openapi.DeregisterOccurrenceResponseObject, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	part, err := h.ops.Deregister(ctx, req.OccurrenceId, ident.Sub.ID)
	if err != nil {
		return nil, err
	}
	return openapi.DeregisterOccurrence200JSONResponse{
		Code: 200, Data: part, Timestamp: time.Now(), Module: module,
	}, nil
}
