package actors

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) GetActorByEmail(ctx context.Context, req openapi.GetActorByEmailRequestObject) (openapi.GetActorByEmailResponseObject, error) {
	err := models.RequireClientOnly(ctx)
	if err != nil {
		return nil, err
	}

	actor, err := h.ops.GetByEmail(ctx, req.ActorEmail, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.GetActorByEmail200JSONResponse{
		Code: 200, Data: actor, Timestamp: time.Now(), Module: module,
	}, nil
}
