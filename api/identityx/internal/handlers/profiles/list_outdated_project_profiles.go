package profiles

import (
	"context"
	"time"

	"IdentityX/internal/openapi"
)

// ListOutdatedProjectProfiles lists the project's actor profiles that
// failed to migrate to the active project schema version.
func (h *Handlers) ListOutdatedProjectProfiles(ctx context.Context, req openapi.ListOutdatedProjectProfilesRequestObject) (openapi.ListOutdatedProjectProfilesResponseObject, error) {
	profiles, err := h.ops.ListOutdatedProfiles(ctx, &req.ProjectId)
	if err != nil {
		return nil, err
	}
	return openapi.ListOutdatedProjectProfiles200JSONResponse{
		Code: 200, Data: profiles, Timestamp: time.Now(), Module: module,
	}, nil
}
