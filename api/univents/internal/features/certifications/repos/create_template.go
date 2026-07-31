package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"
	"univents/models"
)

func (repo *Repo) CreateTemplate(ctx context.Context, input models.CreateCertificationTemplateInput) (*models.CertificationTemplate, error) {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsRepo.CreateTemplate")
	defer span.End()

	result, err := database.Queries(ctx, repo.q).CreateCertificationTemplate(ctx, sqlc.CreateCertificationTemplateParams{
		EditionID:   input.EditionID,
		Kind:        string(input.Kind),
		Name:        input.Name,
		Description: input.Description,
		DesignData:  input.DesignData,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}

	return new(mapCertTemplate(result)), nil
}
