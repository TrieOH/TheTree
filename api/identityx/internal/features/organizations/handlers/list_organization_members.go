package handlers

import (
	"context"

	"IdentityX/internal/openapi"
)

func (h *Handlers) ListOrganizationMembers(ctx context.Context, req openapi.ListOrganizationMembersRequestObject) (openapi.ListOrganizationMembersResponseObject, error) {
	members, err := h.ops.ListMembers(ctx, req.OrganizationId)
	if err != nil {
		return nil, err
	}
	return openapi.ListOrganizationMembers200JSONResponse{Items: members}, nil
}
