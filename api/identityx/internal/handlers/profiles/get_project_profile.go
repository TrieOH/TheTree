package profiles

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
)

func (h *Handlers) GetProjectProfile(ctx context.Context, req openapi.GetProjectProfileRequestObject) (openapi.GetProjectProfileResponseObject, error) {
	profile, err := h.ops.GetProfile(ctx, req.ActorId, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.GetProjectProfile200JSONResponse{
		Code: 200, Data: profile, Timestamp: time.Now(), Module: module,
	}, nil
}
