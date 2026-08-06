package organizations

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) ListOrganizationMembers(ctx context.Context, req openapi.ListOrganizationMembersRequestObject) (openapi.ListOrganizationMembersResponseObject, error) {
	err := models.RequireClientOnly(ctx)
	if err != nil {
		return nil, err
	}

	members, err := h.ops.ListMembers(ctx, req.OrganizationId)
	if err != nil {
		return nil, err
	}
	return openapi.ListOrganizationMembers200JSONResponse{
		Code: 200, Data: &members, Timestamp: time.Now(), Module: module,
	}, nil
}
