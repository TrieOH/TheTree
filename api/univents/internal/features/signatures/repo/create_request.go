package repo

import (
	"context"
	"lib/database"
	"univents/internal/sqlc"
	"univents/models"
)

func (repo *Repo) CreateRequest(ctx context.Context, toCreate *models.SignatureRequest) (*models.SignatureRequest, error) {
	result, err := database.Queries(ctx, repo.q).CreateSignatureRequest(ctx, sqlc.CreateSignatureRequestParams{
		EditionID:       toCreate.EditionID,
		CreatedBy:       toCreate.CreatedBy,
		SignatoryName:   toCreate.SignatoryName,
		SignatoryTitle:  toCreate.SignatoryTitle,
		SignatoryEmail:  toCreate.SignatoryEmail,
		SignatoryUserID: toCreate.SignatoryUserID,
		IdempotencyKey:  toCreate.IdempotencyKey,
		ExpiresAt:       toCreate.ExpiresAt,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapSignatureRequest(result)), nil
}
