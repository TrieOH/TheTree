package handlers

import (
	"context"

	"payssage/internal/openapi"
)

func (h *Handlers) GetOrganizationMemberByEmail(ctx context.Context, req openapi.GetOrganizationMemberByEmailRequestObject) (openapi.GetOrganizationMemberByEmailResponseObject, error) {
	member, err := h.ops.GetMemberByEmail(ctx, req.MemberEmail, req.OrganizationId)
	if err != nil {
		return nil, err
	}
	return openapi.GetOrganizationMemberByEmail200JSONResponse(*member), nil
}
