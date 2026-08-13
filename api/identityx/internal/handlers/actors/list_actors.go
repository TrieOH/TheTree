package actors

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
)

func (h *Handlers) ListActors(ctx context.Context, req openapi.ListActorsRequestObject) (openapi.ListActorsResponseObject, error) {
	actors, err := h.ops.List(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.ListActors200JSONResponse{
		Code: 200, Data: &actors, Timestamp: time.Now(), Module: module,
	}, nil
}
