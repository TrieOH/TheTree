package repos

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (r *Repo) GetTemplateByID(ctx context.Context, id uuid.UUID) (*models.BadgeTemplate, error) {
	row, err := r.q.GetBadgeTemplate(ctx, id)
	if err != nil {
		return nil, r.dbe(err)
	}

	result := mapBadgeTemplate(row)
	return &result, nil
}
