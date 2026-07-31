package queries

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) GetTemplate(ctx context.Context, editionID, templateID uuid.UUID) (*models.BadgeTemplate, error) {
	// Skipping permission checks as requested
	return q.repo.GetTemplateByID(ctx, templateID)
}
