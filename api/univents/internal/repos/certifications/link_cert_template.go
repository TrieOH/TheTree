package certifications

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"

	"github.com/google/uuid"
)

func (repo *Repo) LinkCertTemplate(ctx context.Context, templateID, programID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsRepo.LinkCertTemplate")
	defer span.End()
	return database.Queries(ctx, repo.q).CreateCertificationTemplateProgram(ctx, sqlc.CreateCertificationTemplateProgramParams{
		TemplateID: templateID,
		ProgramID:  programID,
	})
}
