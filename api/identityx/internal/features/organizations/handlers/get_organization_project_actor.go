package handlers

import (
	"context"

	"IdentityX/internal/openapi"
)

func (h *Handlers) GetOrganizationProjectActor(ctx context.Context, req openapi.GetOrganizationProjectActorRequestObject) (openapi.GetOrganizationProjectActorResponseObject, error) {
	actor, err := h.ops.GetActorByID(ctx, req.ActorId, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.GetOrganizationProjectActor200JSONResponse(*actor), nil
}
