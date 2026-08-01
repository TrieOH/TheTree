package certifications

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"

	"github.com/google/uuid"
)

func (repo *Repo) UnlinkCertTemplate(ctx context.Context, templateID, programID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsRepo.UnlinkCertTemplate")
	defer span.End()
	return database.Queries(ctx, repo.q).DeleteCertificationTemplatePrograms(ctx, sqlc.DeleteCertificationTemplateProgramsParams{
		TemplateID: templateID,
		ProgramID:  programID,
	})
}
