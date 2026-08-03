package profiles

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
	"IdentityX/models"
)

func (h *Handlers) UpsertProjectProfile(ctx context.Context, req openapi.UpsertProjectProfileRequestObject) (openapi.UpsertProjectProfileResponseObject, error) {
	profile, err := h.ops.UpsertProfile(ctx, models.UpsertProfileInput{
		ActorID: req.ActorId,
		Profile: req.Body.Profile,
	}, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.UpsertProjectProfile200JSONResponse{
		Code: 200, Data: profile, Timestamp: time.Now(), Module: module,
	}, nil
}
