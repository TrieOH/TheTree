package repo

import (
	"context"

	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) ResponsesCount(ctx context.Context, formID uuid.UUID) (int, error) {
	count, err := database.Queries(ctx, repo.q).CountFormResponses(ctx, formID)
	if err != nil {
		return 0, repo.dbe(err)
	}
	return int(count), nil
}
