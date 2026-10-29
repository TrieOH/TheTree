package handlers

import (
	"context"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) RemoveOrganizationMember(ctx context.Context, req openapi.RemoveOrganizationMemberRequestObject) (openapi.RemoveOrganizationMemberResponseObject, error) {
	err := h.ops.RemoveMember(ctx, models.RemoveOrganizationMemberInput{
		ActorEmail:     req.Body.ActorEmail,
		OrganizationID: req.OrganizationId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.RemoveOrganizationMember204Response{}, nil
}
