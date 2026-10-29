package handlers

import (
	"context"

	"IdentityX/internal/openapi"
)

func (h *Handlers) GetActorByEmail(ctx context.Context, req openapi.GetActorByEmailRequestObject) (openapi.GetActorByEmailResponseObject, error) {
	actor, err := h.ops.GetByEmail(ctx, req.ActorEmail, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.GetActorByEmail200JSONResponse(*actor), nil
}
