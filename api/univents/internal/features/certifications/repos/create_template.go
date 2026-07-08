package repos

import (
	"context"
	"lib/database"
	"univents/contracts"
	"univents/internal/database/sqlc"
)

func (repo *repo) CreateTemplate(ctx context.Context, input contracts.CreateCertificationTemplateInput) (*contracts.CertificationTemplate, error) {
	ctx, span := database.Span(ctx, repo.tracer, "CreateTemplate")
	defer span.End()
	row, err := database.Queries(ctx, repo.q).CreateCertificationTemplate(ctx, sqlc.CreateCertificationTemplateParams{
		EditionID: input.EditionID,
		Title:     input.Title,
		Data:      input.Data,
		Url:       input.URL,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapCertificationTemplate(row)), nil
}
