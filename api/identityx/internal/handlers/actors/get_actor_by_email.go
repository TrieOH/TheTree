package actors

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
)

func (h *Handlers) GetActorByEmail(ctx context.Context, req openapi.GetActorByEmailRequestObject) (openapi.GetActorByEmailResponseObject, error) {
	actor, err := h.ops.GetByEmail(ctx, req.ActorEmail, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.GetActorByEmail200JSONResponse{
		Code: 200, Data: actor, Timestamp: time.Now(), Module: module,
	}, nil
}
