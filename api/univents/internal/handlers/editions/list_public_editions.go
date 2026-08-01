package editions

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) ListPublicEditions(ctx context.Context, req openapi.ListPublicEditionsRequestObject) (openapi.ListPublicEditionsResponseObject, error) {
	editions, err := h.ops.ListPublic(ctx, req.EventId)
	if err != nil {
		return nil, err
	}
	return openapi.ListPublicEditions200JSONResponse{
		Code: 200, Data: &editions, Timestamp: time.Now(), Module: module,
	}, nil
}
