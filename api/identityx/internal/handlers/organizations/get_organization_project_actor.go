package organizations

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) GetOrganizationProjectActor(ctx context.Context, req openapi.GetOrganizationProjectActorRequestObject) (openapi.GetOrganizationProjectActorResponseObject, error) {
	err := models.RequireClientOnly(ctx)
	if err != nil {
		return nil, err
	}

	actor, err := h.ops.GetActorByID(ctx, req.ActorId, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.GetOrganizationProjectActor200JSONResponse{
		Code: 200, Data: actor, Timestamp: time.Now(), Module: module,
	}, nil
}
