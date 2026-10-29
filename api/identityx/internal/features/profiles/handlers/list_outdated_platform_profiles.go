package handlers

import (
	"context"

	"IdentityX/internal/openapi"
)

// ListOutdatedPlatformProfiles lists platform-scoped actor profiles that
// failed to migrate to the active platform schema version.
func (h *Handlers) ListOutdatedPlatformProfiles(ctx context.Context, _ openapi.ListOutdatedPlatformProfilesRequestObject) (openapi.ListOutdatedPlatformProfilesResponseObject, error) {
	profiles, err := h.ops.ListOutdatedProfiles(ctx, nil)
	if err != nil {
		return nil, err
	}
	return openapi.ListOutdatedPlatformProfiles200JSONResponse{Items: profiles}, nil
}
