package authn

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/telemetry"
)

func (repo *Repo) GetByProviderAndSubject(ctx context.Context, provider, subject string) (*models.ActorExternalIdentities, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetByProviderAndSubject")
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
