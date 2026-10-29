package profiles

import (
	"context"

	"IdentityX/internal/features/profiles/service"
	"IdentityX/internal/openapi"
	"IdentityX/models"
)

type Handlers struct {
	ops *service.Operations
}

func NewHandlers(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) GetProfileByHandle(ctx context.Context, req openapi.GetProfileByHandleRequestObject) (openapi.GetProfileByHandleResponseObject, error) {
	profile, err := h.ops.GetProfileByHandle(ctx, req.Handle)
	if err != nil {
		return nil, err
	}
	return openapi.GetProfileByHandle200JSONResponse(*profile), nil
}

func (h *Handlers) GetPlatformProfile(ctx context.Context, req openapi.GetPlatformProfileRequestObject) (openapi.GetPlatformProfileResponseObject, error) {
	profile, err := h.ops.GetPlatformProfile(ctx, req.ActorId)
	if err != nil {
		return nil, err
	}
	return openapi.GetPlatformProfile200JSONResponse(*profile), nil
}

func (h *Handlers) GetProjectProfile(ctx context.Context, req openapi.GetProjectProfileRequestObject) (openapi.GetProjectProfileResponseObject, error) {
	profile, err := h.ops.GetProfile(ctx, req.ActorId, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.GetProjectProfile200JSONResponse(*profile), nil
}

// ListOutdatedPlatformProfiles lists platform-scoped actor profiles that
// failed to migrate to the active platform schema version.
func (h *Handlers) ListOutdatedPlatformProfiles(ctx context.Context, _ openapi.ListOutdatedPlatformProfilesRequestObject) (openapi.ListOutdatedPlatformProfilesResponseObject, error) {
	profiles, err := h.ops.ListOutdatedProfiles(ctx, nil)
	if err != nil {
		return nil, err
	}
	return openapi.ListOutdatedPlatformProfiles200JSONResponse{Items: profiles}, nil
}

// ListOutdatedProjectProfiles lists the project's actor profiles that
// failed to migrate to the active project schema version.
func (h *Handlers) ListOutdatedProjectProfiles(ctx context.Context, req openapi.ListOutdatedProjectProfilesRequestObject) (openapi.ListOutdatedProjectProfilesResponseObject, error) {
	profiles, err := h.ops.ListOutdatedProfiles(ctx, &req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.ListOutdatedProjectProfiles200JSONResponse{Items: profiles}, nil
}

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
	return openapi.UpsertPlatformProfile200JSONResponse(*profile), nil
}

func (h *Handlers) UpsertProjectProfile(ctx context.Context, req openapi.UpsertProjectProfileRequestObject) (openapi.UpsertProjectProfileResponseObject, error) {
	profile, err := h.ops.UpsertProfile(ctx, models.UpsertProfileInput{
		ActorID: req.ActorId,
		Handle:  req.Body.Handle,
		PfpURL:  req.Body.PfpUrl,
		Profile: req.Body.Profile,
	}, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.UpsertProjectProfile200JSONResponse(*profile), nil
}
