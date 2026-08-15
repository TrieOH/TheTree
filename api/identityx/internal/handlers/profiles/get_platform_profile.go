package profiles

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
)

func (h *Handlers) GetPlatformProfile(ctx context.Context, req openapi.GetPlatformProfileRequestObject) (openapi.GetPlatformProfileResponseObject, error) {
	profile, err := h.ops.GetPlatformProfile(ctx, req.ActorId)
	if err != nil {
		return nil, err
	}
	return openapi.GetPlatformProfile200JSONResponse{
		Code: 200, Data: profile, Timestamp: time.Now(), Module: module,
	}, nil
}
