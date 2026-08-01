package editions

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) ListDraftEditions(ctx context.Context, req openapi.ListDraftEditionsRequestObject) (openapi.ListDraftEditionsResponseObject, error) {
	editions, err := h.ops.ListDraft(ctx, req.EventId)
	if err != nil {
		return nil, err
	}
	return openapi.ListDraftEditions200JSONResponse{
		Code: 200, Data: &editions, Timestamp: time.Now(), Module: module,
	}, nil
}
