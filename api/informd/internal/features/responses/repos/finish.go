package repos

import (
	"context"

	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) Finish(ctx context.Context, id uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "ResponseRepo.Finish")
	defer span.End()
	err := database.Queries(ctx, repo.q).FinishResponse(ctx, id)
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
