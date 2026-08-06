package organizations

import (
	"context"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) AddOrganizationMember(ctx context.Context, req openapi.AddOrganizationMemberRequestObject) (openapi.AddOrganizationMemberResponseObject, error) {
	err := models.RequireClientOnly(ctx)
	if err != nil {
		return nil, err
	}

	err = h.ops.AddMember(ctx, models.AddOrganizationMemberInput{
		ActorEmail:     req.Body.ActorEmail,
		Role:           req.Body.Role,
		OrganizationID: req.OrganizationId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.AddOrganizationMember201Response{}, nil
}
