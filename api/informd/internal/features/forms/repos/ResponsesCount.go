package repos

import (
	"context"

	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) ResponsesCount(ctx context.Context, formID uuid.UUID) (int, error) {
	ctx, span := telemetry.StartSpan(ctx, "FormRepo.ResponsesCount")
	defer span.End()
	count, err := database.Queries(ctx, repo.q).CountFormResponses(ctx, formID)
	if err != nil {
		return 0, repo.dbe(err)
	}
	return int(count), nil
}
