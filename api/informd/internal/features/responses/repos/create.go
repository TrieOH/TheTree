package repos

import (
	"Informd/internal/database/sqlc"
	"Informd/models"
	"context"
	"lib/database"
)

func (repo *repo) Create(ctx context.Context, toCreate models.Response) (*models.Response, error) {
	ctx, span := repo.tracer.Start(ctx, "ResponseRepo.Create")
	defer span.End()
	row, err := database.Queries(ctx, repo.q).CreateResponse(ctx, sqlc.CreateResponseParams{
		FormID:      toCreate.FormID,
		InviteID:    toCreate.InviteID,
		ResponderID: toCreate.ResponderID,
		Email:       toCreate.Email,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapResponse(row)), nil
}
