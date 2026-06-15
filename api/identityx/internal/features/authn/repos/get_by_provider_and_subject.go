package repos

import (
	"IdentityX/internal/database/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"
)

func (repo *repo) GetByProviderAndSubject(ctx context.Context, provider, subject string) (*models.ActorExternalIdentities, error) {
	ctx, span := database.Span(ctx, repo.tracer, "GetByProviderAndSubject")
	defer span.End()
	row, err := database.Queries(ctx, repo.q).GetExternalIdentityByProviderAndSubject(ctx, sqlc.GetExternalIdentityByProviderAndSubjectParams{
		Provider: provider,
		Subject:  subject,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapExternalIdentity(row)), nil
}
