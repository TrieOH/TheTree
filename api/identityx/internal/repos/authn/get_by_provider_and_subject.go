package authn

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) GetByProviderAndSubject(ctx context.Context, provider, subject string, projectID *uuid.UUID) (*models.ActorExternalIdentities, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetByProviderAndSubject")
	defer span.End()

	row, err := database.Queries(ctx, repo.q).GetExternalIdentityByProviderAndSubject(ctx, sqlc.GetExternalIdentityByProviderAndSubjectParams{
		Provider:  provider,
		Subject:   subject,
		ProjectID: projectID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapExternalIdentity(row)), nil
}
