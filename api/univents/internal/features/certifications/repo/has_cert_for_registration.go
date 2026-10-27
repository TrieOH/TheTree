package repo

import (
	"context"
	"lib/database"
	"univents/internal/sqlc"

	"github.com/google/uuid"
)

func (repo *Repo) HasCertForRegistration(ctx context.Context, registrationID uuid.UUID, templateID *uuid.UUID) (bool, error) {
	return database.Queries(ctx, repo.q).HasCertForRegistration(ctx, sqlc.HasCertForRegistrationParams{
		RegistrationID: registrationID,
		TemplateID:     templateID,
	})
}
