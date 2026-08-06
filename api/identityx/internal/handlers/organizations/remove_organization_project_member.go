package organizations

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) RemoveOrganizationProjectMember(ctx context.Context, req openapi.RemoveOrganizationProjectMemberRequestObject) (openapi.RemoveOrganizationProjectMemberResponseObject, error) {
	err := models.RequireClientOnly(ctx)
	if err != nil {
		return nil, err
	}

	err = h.ops.RemoveProjectMember(ctx, models.RemoveOrgProjectMemberInput{
		ActorEmail:     req.Body.ActorEmail,
		OrganizationID: req.OrganizationId,
		ProjectID:      req.ProjectId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.RemoveOrganizationProjectMember200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}
