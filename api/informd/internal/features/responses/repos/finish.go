package repos

import (
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) Finish(ctx context.Context, id uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "ResponseRepo.Finish")
	defer span.End()
	err := database.Queries(ctx, repo.q).FinishResponse(ctx, id)
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
