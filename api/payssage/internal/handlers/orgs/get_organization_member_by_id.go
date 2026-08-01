package orgs

import (
	"context"
	"time"

	"payssage/internal/openapi"
)

func (h *Handlers) GetOrganizationMemberByID(ctx context.Context, req openapi.GetOrganizationMemberByIDRequestObject) (openapi.GetOrganizationMemberByIDResponseObject, error) {
	member, err := h.ops.GetMemberByID(ctx, req.MemberId, req.OrganizationId)
	if err != nil {
		return nil, err
	}
	return openapi.GetOrganizationMemberByID200JSONResponse{
		Code: 200, Data: member, Timestamp: time.Now(), Module: module,
	}, nil
}
