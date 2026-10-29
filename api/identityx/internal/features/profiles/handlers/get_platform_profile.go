package handlers

import (
	"context"

	"IdentityX/internal/openapi"
)

func (h *Handlers) GetPlatformProfile(ctx context.Context, req openapi.GetPlatformProfileRequestObject) (openapi.GetPlatformProfileResponseObject, error) {
	profile, err := h.ops.GetPlatformProfile(ctx, req.ActorId)
	if err != nil {
		return nil, err
	}
	return openapi.GetPlatformProfile200JSONResponse(*profile), nil
}
