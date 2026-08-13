package actors

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
)

func (h *Handlers) GetActor(ctx context.Context, req openapi.GetActorRequestObject) (openapi.GetActorResponseObject, error) {
	actor, err := h.ops.GetByID(ctx, req.ActorId, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.GetActor200JSONResponse{
		Code: 200, Data: actor, Timestamp: time.Now(), Module: module,
	}, nil
}
