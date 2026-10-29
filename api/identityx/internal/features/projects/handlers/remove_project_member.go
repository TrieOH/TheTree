package handlers

import (
	"context"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) RemoveProjectMember(ctx context.Context, req openapi.RemoveProjectMemberRequestObject) (openapi.RemoveProjectMemberResponseObject, error) {
	err := h.ops.RemoveMember(ctx, models.RemoveProjectMemberInput{
		ActorEmail: req.Body.ActorEmail,
		ProjectID:  req.ProjectId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.RemoveProjectMember204Response{}, nil
}
