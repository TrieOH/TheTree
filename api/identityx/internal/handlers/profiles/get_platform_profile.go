package profiles

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) GetPlatformProfile(ctx context.Context, req openapi.GetPlatformProfileRequestObject) (openapi.GetPlatformProfileResponseObject, error) {
	err := models.RequireClientOnly(ctx)
	if err != nil {
		return nil, err
	}

	profile, err := h.ops.GetPlatformProfile(ctx, req.ActorId)
	if err != nil {
		return nil, err
	}
	return openapi.GetPlatformProfile200JSONResponse{
		Code: 200, Data: profile, Timestamp: time.Now(), Module: module,
	}, nil
}
