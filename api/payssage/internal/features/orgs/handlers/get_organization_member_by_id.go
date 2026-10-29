package handlers

import (
	"context"

	"payssage/internal/openapi"
)

func (h *Handlers) GetOrganizationMemberByID(ctx context.Context, req openapi.GetOrganizationMemberByIDRequestObject) (openapi.GetOrganizationMemberByIDResponseObject, error) {
	member, err := h.ops.GetMemberByID(ctx, req.MemberId, req.OrganizationId)
	if err != nil {
		return nil, err
	}
	return openapi.GetOrganizationMemberByID200JSONResponse(*member), nil
}
