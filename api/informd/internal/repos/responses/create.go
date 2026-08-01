package responses

import (
	"Informd/internal/sqlc"
	"context"

	"Informd/models"
	"lib/database"
	"lib/telemetry"
)

func (repo *Repo) Create(ctx context.Context, toCreate models.Response) (*models.Response, error) {
	ctx, span := telemetry.StartSpan(ctx, "ResponseRepo.Create")
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
