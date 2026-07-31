package repos

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (r *Repo) ListTemplatesByEdition(ctx context.Context, editionID uuid.UUID) ([]models.BadgeTemplate, error) {
	rows, err := r.q.ListBadgeTemplatesByEdition(ctx, editionID)
	if err != nil {
		return nil, r.dbe(err)
	}

	result := make([]models.BadgeTemplate, len(rows))
	for i, row := range rows {
		result[i] = mapBadgeTemplate(row)
	}
	return result, nil
}
