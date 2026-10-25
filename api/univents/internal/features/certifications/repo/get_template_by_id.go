package repo

import (
	"context"
	"lib/database"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetTemplateByID(ctx context.Context, id uuid.UUID) (*models.CertificationTemplate, error) {
	result, err := database.Queries(ctx, repo.q).GetCertificationTemplateByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapCertTemplate(result)), nil
}
