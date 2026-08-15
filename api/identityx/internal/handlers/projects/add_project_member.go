package projects

import (
	"context"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) AddProjectMember(ctx context.Context, req openapi.AddProjectMemberRequestObject) (openapi.AddProjectMemberResponseObject, error) {
	err := h.ops.AddMember(ctx, models.AddProjectMemberInput{
		ActorEmail: req.Body.ActorEmail,
		Role:       req.Body.Role,
		ProjectID:  req.ProjectId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.AddProjectMember201Response{}, nil
}
