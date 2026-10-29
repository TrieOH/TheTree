package handlers

import (
	"context"

	"IdentityX/internal/openapi"
)

func (h *Handlers) ListProjectMembers(ctx context.Context, req openapi.ListProjectMembersRequestObject) (openapi.ListProjectMembersResponseObject, error) {
	members, err := h.ops.ListMembers(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.ListProjectMembers200JSONResponse{Items: members}, nil
}
