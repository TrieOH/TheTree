package editions

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) ListPastEditions(ctx context.Context, req openapi.ListPastEditionsRequestObject) (openapi.ListPastEditionsResponseObject, error) {
	editions, err := h.ops.GetPast(ctx, req.EventId)
	if err != nil {
		return nil, err
	}
	return openapi.ListPastEditions200JSONResponse{
		Code: 200, Data: &editions, Timestamp: time.Now(), Module: module,
	}, nil
}
