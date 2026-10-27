package repo

import (
	"context"
	"lib/database"
	"univents/internal/sqlc"
	"univents/models"
)

func (repo *Repo) UpdateTemplate(ctx context.Context, input models.UpdateCertificationTemplateInput) (*models.CertificationTemplate, error) {
	result, err := database.Queries(ctx, repo.q).UpdateCertificationTemplate(ctx, sqlc.UpdateCertificationTemplateParams{
		ID:          input.TemplateID,
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
