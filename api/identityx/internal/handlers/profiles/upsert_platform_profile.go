package profiles

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) UpsertPlatformProfile(ctx context.Context, req openapi.UpsertPlatformProfileRequestObject) (openapi.UpsertPlatformProfileResponseObject, error) {
	profile, err := h.ops.UpsertPlatformProfile(ctx, models.UpsertProfileInput{
		ActorID: req.ActorId,
		Handle:  req.Body.Handle,
		PfpURL:  req.Body.PfpUrl,
		Profile: req.Body.Profile,
	})
	if err != nil {
		return nil, err
	}
	return openapi.UpsertPlatformProfile200JSONResponse{
		Code: 200, Data: profile, Timestamp: time.Now(), Module: module,
	}, nil
}
