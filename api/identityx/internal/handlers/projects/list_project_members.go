package projects

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) ListProjectMembers(ctx context.Context, req openapi.ListProjectMembersRequestObject) (openapi.ListProjectMembersResponseObject, error) {
	err := models.RequireClientOnly(ctx)
	if err != nil {
		return nil, err
	}

	members, err := h.ops.ListMembers(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.ListProjectMembers200JSONResponse{
		Code: 200, Data: &members, Timestamp: time.Now(), Module: module,
	}, nil
}
