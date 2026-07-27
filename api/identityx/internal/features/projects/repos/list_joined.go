package repos

import (
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
	"lib/telemetry"
)

func (repo *Repo) ListJoined(ctx context.Context, userID uuid.UUID) ([]models.Project, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListJoined")
	defer span.End()

	sqlcProjects, err := database.Queries(ctx, repo.q).ListJoinedProjects(ctx, userID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcProjects, mapProject), nil
}
