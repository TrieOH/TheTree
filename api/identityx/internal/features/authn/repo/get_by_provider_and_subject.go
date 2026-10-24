package repo

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) GetByProviderAndSubject(ctx context.Context, provider, subject string, projectID *uuid.UUID) (*models.ActorExternalIdentities, error) {
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
