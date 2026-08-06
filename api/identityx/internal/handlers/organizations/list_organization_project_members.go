package organizations

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) ListOrganizationProjectMembers(ctx context.Context, req openapi.ListOrganizationProjectMembersRequestObject) (openapi.ListOrganizationProjectMembersResponseObject, error) {
	err := models.RequireClientOnly(ctx)
	if err != nil {
		return nil, err
	}

	members, err := h.ops.ListOrgProjectMembers(ctx, req.OrganizationId, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.ListOrganizationProjectMembers200JSONResponse{
		Code: 200, Data: &members, Timestamp: time.Now(), Module: module,
	}, nil
}
