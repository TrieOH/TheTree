package organizations

import (
	"context"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) AddOrganizationProjectMember(ctx context.Context, req openapi.AddOrganizationProjectMemberRequestObject) (openapi.AddOrganizationProjectMemberResponseObject, error) {
	err := h.ops.AddProjectMember(ctx, models.AddOrgProjectMemberInput{
		ActorEmail:     req.Body.ActorEmail,
		Role:           req.Body.Role,
		OrganizationID: req.OrganizationId,
		ProjectID:      req.ProjectId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.AddOrganizationProjectMember201Response{}, nil
}
