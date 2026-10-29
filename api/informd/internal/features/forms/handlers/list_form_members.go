package handlers

import (
	"context"

	"Informd/internal/openapi"
)

func (h *Handlers) ListFormMembers(ctx context.Context, req openapi.ListFormMembersRequestObject) (openapi.ListFormMembersResponseObject, error) {
	members, err := h.ops.ListMembers(ctx, req.FormId)
	if err != nil {
		return nil, err
	}
	return openapi.ListFormMembers200JSONResponse{Items: members}, nil
}
