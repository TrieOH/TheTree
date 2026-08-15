package projects

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
)

func (h *Handlers) ListProjectMembers(ctx context.Context, req openapi.ListProjectMembersRequestObject) (openapi.ListProjectMembersResponseObject, error) {
	members, err := h.ops.ListMembers(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.ListProjectMembers200JSONResponse{
		Code: 200, Data: &members, Timestamp: time.Now(), Module: module,
	}, nil
}
