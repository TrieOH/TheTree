package orgs

import (
	"context"
	"time"

	"payssage/internal/openapi"
)

func (h *Handlers) GetOrganizationMemberByEmail(ctx context.Context, req openapi.GetOrganizationMemberByEmailRequestObject) (openapi.GetOrganizationMemberByEmailResponseObject, error) {
	member, err := h.ops.GetMemberByEmail(ctx, req.MemberEmail, req.OrganizationId)
	if err != nil {
		return nil, err
	}
	return openapi.GetOrganizationMemberByEmail200JSONResponse{
		Code: 200, Data: member, Timestamp: time.Now(), Module: module,
	}, nil
}
