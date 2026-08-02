package profiles

import (
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/xslices"

	"lib/telemetry"

	"github.com/google/uuid"
)

// ListOutdated returns profiles flagged as outdated. A nil projectID means
// the platform scope (actors with no project).
func (r *Repo) ListOutdated(ctx context.Context, projectID *uuid.UUID) ([]models.ActorProfile, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListOutdated")
	defer span.End()

	results, err := database.Queries(ctx, r.q).ListOutdatedProfiles(ctx, projectID)
	if err != nil {
		return nil, r.dbe(err)
	}
	return xslices.MapSlice(results, mapActorProfile), nil
}
