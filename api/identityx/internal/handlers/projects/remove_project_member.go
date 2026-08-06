package projects

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) RemoveProjectMember(ctx context.Context, req openapi.RemoveProjectMemberRequestObject) (openapi.RemoveProjectMemberResponseObject, error) {
	err := models.RequireClientOnly(ctx)
	if err != nil {
		return nil, err
	}

	err = h.ops.RemoveMember(ctx, models.RemoveProjectMemberInput{
		ActorEmail: req.Body.ActorEmail,
		ProjectID:  req.ProjectId,
	})
	if err != nil {
		return nil, err
	}
	return openapi.RemoveProjectMember200JSONResponse{
		Code: 200, Timestamp: time.Now(), Module: module,
	}, nil
}
