package handlers

import (
	"context"

	"IdentityX/internal/openapi"
)

func (h *Handlers) ListOrganizationProjectMembers(ctx context.Context, req openapi.ListOrganizationProjectMembersRequestObject) (openapi.ListOrganizationProjectMembersResponseObject, error) {
	members, err := h.ops.ListOrgProjectMembers(ctx, req.OrganizationId, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.ListOrganizationProjectMembers200JSONResponse{Items: members}, nil
}
