package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"
	"univents/models"
)

func (repo *Repo) Create(ctx context.Context, toCreate *models.Signature) (*models.Signature, error) {
	ctx, span := telemetry.StartSpan(ctx, "SignaturesRepo.Create")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).CreateSignature(ctx, sqlc.CreateSignatureParams{
		EditionID:       toCreate.EditionID,
		CreatedBy:       toCreate.CreatedBy,
		SignatoryName:   toCreate.SignatoryName,
		SignatoryTitle:  toCreate.SignatoryTitle,
		SignatoryEmail:  toCreate.SignatoryEmail,
		SignatoryUserID: toCreate.SignatoryUserID,
		ImageUrl:        toCreate.ImageURL,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapSignature(result)), nil
}
