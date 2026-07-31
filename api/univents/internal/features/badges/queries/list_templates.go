package queries

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) ListTemplates(ctx context.Context, editionID uuid.UUID) ([]models.BadgeTemplate, error) {
	// Skipping permission checks as requested
	return q.repo.ListTemplatesByEdition(ctx, editionID)
}
