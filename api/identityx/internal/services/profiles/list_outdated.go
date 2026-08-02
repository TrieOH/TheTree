package profiles

import (
	"IdentityX/models"
	"context"

	"lib/telemetry"

	"github.com/google/uuid"
)

// ListOutdatedProfiles lists profiles flagged as outdated (they failed to
// migrate to the active schema version). A nil projectID means the platform
// scope; the project variant requires a project admin.
func (o *Operations) ListOutdatedProfiles(ctx context.Context, projectID *uuid.UUID) ([]models.ActorProfile, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListOutdatedProfiles")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	if projectID != nil {
		err = o.authz.CheckProject(ctx, ident.Sub.ID, *projectID, nil, models.ProjectRoleAdmin)
		if err != nil {
			return nil, err
		}
	}

	return o.profiles.ListOutdated(ctx, projectID)
}
