package certifications

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"

	"github.com/google/uuid"
)

func (repo *Repo) HasCertForRegistration(ctx context.Context, registrationID uuid.UUID, templateID *uuid.UUID) (bool, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsRepo.HasCertForRegistration")
	defer span.End()
	return database.Queries(ctx, repo.q).HasCertForRegistration(ctx, sqlc.HasCertForRegistrationParams{
		RegistrationID: registrationID,
		TemplateID:     templateID,
	})
}
