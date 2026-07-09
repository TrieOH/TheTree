package repos

import (
	"context"
	"errors"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *repo) AddEdition(ctx context.Context, eventID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.AddEdition")
	defer span.End()

	affectedRows, err := database.Queries(ctx, repo.q).AddEdition(ctx, eventID)
	if err != nil {
		return repo.dbe(err)
	}

	if affectedRows == 0 {
		return errors.New("could not add edition")
	}

	return nil
}
