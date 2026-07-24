package repos

import (
	"context"
	"lib/database"
	"univents/internal/database/sqlc"
	"univents/models"
)

func (repo *repo) Add(ctx context.Context, toAdd models.Signature) (*models.Signature, error) {
	ctx, span := database.Span(ctx, repo.tracer, "Add")
	defer span.End()
	row, err := database.Queries(ctx, repo.q).AddSignatureToEdition(ctx, sqlc.AddSignatureToEditionParams{
		EditionID: toAdd.EditionID,
		Title:     toAdd.Title,
		Url:       toAdd.URL,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapSignature(row)), nil
}
