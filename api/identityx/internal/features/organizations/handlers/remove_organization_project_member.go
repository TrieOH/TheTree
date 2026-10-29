package handlers

import (
	"context"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) RemoveOrganizationProjectMember(ctx context.Context, req openapi.RemoveOrganizationProjectMemberRequestObject) (openapi.RemoveOrganizationProjectMemberResponseObject, error) {
	err := h.ops.RemoveProjectMember(ctx, models.RemoveOrgProjectMemberInput{
		ActorEmail:     req.Body.ActorEmail,
		OrganizationID: req.OrganizationId,
		ProjectID:      req.ProjectId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.RemoveOrganizationProjectMember204Response{}, nil
}
