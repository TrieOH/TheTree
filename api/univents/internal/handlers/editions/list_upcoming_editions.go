package editions

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) ListUpcomingEditions(ctx context.Context, req openapi.ListUpcomingEditionsRequestObject) (openapi.ListUpcomingEditionsResponseObject, error) {
	editions, err := h.ops.GetUpcoming(ctx, req.EventId)
	if err != nil {
		return nil, err
	}
	return openapi.ListUpcomingEditions200JSONResponse{
		Code: 200, Data: &editions, Timestamp: time.Now(), Module: module,
	}, nil
}
