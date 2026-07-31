package commands

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (c *Commands) CreateTemplate(ctx context.Context, input models.CreateBadgeTemplateInput) (*models.BadgeTemplate, error) {
	// Skipping permission checks as requested
	template := &models.BadgeTemplate{
		ID:           uuid.New(),
		EditionID:    input.EditionID,
		TicketTypeID: input.TicketTypeID,
		Name:         input.Name,
		DesignData:   input.DesignData,
	}

	return c.repo.CreateTemplate(ctx, template)
}
