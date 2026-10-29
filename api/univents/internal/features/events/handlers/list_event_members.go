package handlers

import (
	"context"

	"univents/internal/openapi"
)

func (h *Handlers) ListEventMembers(ctx context.Context, req openapi.ListEventMembersRequestObject) (openapi.ListEventMembersResponseObject, error) {
	members, err := h.ops.ListMembers(ctx, req.EventId)
	if err != nil {
		return nil, err
	}
	return openapi.ListEventMembers200JSONResponse{Items: members}, nil
}
