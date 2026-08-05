package certifications

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetTemplateByID(ctx context.Context, id uuid.UUID) (*models.CertificationTemplate, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsRepo.GetTemplateByID")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).GetCertificationTemplateByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapCertTemplate(result)), nil
}
