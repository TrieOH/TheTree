package repos

import (
	"context"
	"lib/database"
	"univents/contracts"
	"univents/internal/database/sqlc"
)

func (repo *repo) Add(ctx context.Context, toAdd contracts.Signature) (*contracts.Signature, error) {
	ctx, span := database.Span(ctx, repo.tracer, "Add")
	defer span.End()
	row, err := database.Queries(ctx, repo.q).AddSignatureToEdition(ctx, sqlc.AddSignatureToEditionParams{
		ID:        toAdd.ID,
		EditionID: toAdd.EditionID,
		Title:     toAdd.Title,
		Url:       toAdd.URL,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapSignature(row)), nil
}
