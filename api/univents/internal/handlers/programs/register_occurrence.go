package programs

import (
	"context"
	"time"

	idx "sdk/identityx"

	"univents/internal/openapi"
)

func (h *Handlers) RegisterOccurrence(ctx context.Context, req openapi.RegisterOccurrenceRequestObject) (openapi.RegisterOccurrenceResponseObject, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	part, err := h.ops.Register(ctx, req.OccurrenceId, ident.Sub.ID)
	if err != nil {
		return nil, err
	}
	return openapi.RegisterOccurrence201JSONResponse{
		Code: 201, Data: part, Timestamp: time.Now(), Module: module,
	}, nil
}
