package organizations

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
)

func (h *Handlers) GetOrganizationProjectActor(ctx context.Context, req openapi.GetOrganizationProjectActorRequestObject) (openapi.GetOrganizationProjectActorResponseObject, error) {
	actor, err := h.ops.GetActorByID(ctx, req.ActorId, req.OrganizationId, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.GetOrganizationProjectActor200JSONResponse{
		Code: 200, Data: actor, Timestamp: time.Now(), Module: module,
	}, nil
}
