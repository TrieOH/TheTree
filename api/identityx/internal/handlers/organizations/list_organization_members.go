package organizations

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
)

func (h *Handlers) ListOrganizationMembers(ctx context.Context, req openapi.ListOrganizationMembersRequestObject) (openapi.ListOrganizationMembersResponseObject, error) {
	members, err := h.ops.ListMembers(ctx, req.OrganizationId)
	if err != nil {
		return nil, err
	}
	return openapi.ListOrganizationMembers200JSONResponse{
		Code: 200, Data: &members, Timestamp: time.Now(), Module: module,
	}, nil
}
