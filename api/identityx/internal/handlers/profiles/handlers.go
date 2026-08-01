// Package profiles implements the StrictServerInterface methods for the
// profiles feature.
package profiles

import (
	"context"
	"encoding/json"
	"time"

	"IdentityX/internal/openapi"
	"IdentityX/internal/services"
	"IdentityX/models"
)

const module = "IdentityX"

type Handlers struct {
	ops *services.Profiles
}

func New(ops *services.Profiles) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) GetPlatformProfile(ctx context.Context, req openapi.GetPlatformProfileRequestObject) (openapi.GetPlatformProfileResponseObject, error) {
	profile, err := h.ops.GetPlatformProfile(ctx, req.ActorId)
	if err != nil {
		return nil, err
	}
	return openapi.GetPlatformProfile200JSONResponse{
		Code: 200, Data: profile, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) UpsertPlatformProfile(ctx context.Context, req openapi.UpsertPlatformProfileRequestObject) (openapi.UpsertPlatformProfileResponseObject, error) {
	profile, err := h.ops.UpsertPlatformProfile(ctx, models.UpsertProfileInput{
		ActorID: req.ActorId,
		Profile: json.RawMessage(req.Body.Profile),
	})
	if err != nil {
		return nil, err
	}
	return openapi.UpsertPlatformProfile200JSONResponse{
		Code: 200, Data: profile, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetProjectProfile(ctx context.Context, req openapi.GetProjectProfileRequestObject) (openapi.GetProjectProfileResponseObject, error) {
	profile, err := h.ops.GetProfile(ctx, req.ActorId, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.GetProjectProfile200JSONResponse{
		Code: 200, Data: profile, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) UpsertProjectProfile(ctx context.Context, req openapi.UpsertProjectProfileRequestObject) (openapi.UpsertProjectProfileResponseObject, error) {
	profile, err := h.ops.UpsertProfile(ctx, models.UpsertProfileInput{
		ActorID: req.ActorId,
		Profile: json.RawMessage(req.Body.Profile),
	}, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.UpsertProjectProfile200JSONResponse{
		Code: 200, Data: profile, Timestamp: time.Now(), Module: module,
	}, nil
}
