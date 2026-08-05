package events

import (
	"context"
	"time"

	"univents/internal/openapi"
)

func (h *Handlers) ListEventMembers(ctx context.Context, req openapi.ListEventMembersRequestObject) (openapi.ListEventMembersResponseObject, error) {
	members, err := h.ops.ListMembers(ctx, req.EventId)
	if err != nil {
		return nil, err
	}
	return openapi.ListEventMembers200JSONResponse{
		Code: 200, Data: &members, Timestamp: time.Now(), Module: module,
	}, nil
}
