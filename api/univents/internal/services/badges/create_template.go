package badges

import (
	"context"
	"lib/telemetry"
	"univents/models"

	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (o *Operations) CreateTemplate(ctx context.Context, input models.CreateBadgeTemplateInput) (*models.BadgeTemplate, error) {
	ctx, span := telemetry.StartSpan(ctx, "BadgeTemplatesCommands.CreateTemplate")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	edition, err := o.editions.GetByID(ctx, input.EditionID)
	if err != nil {
		return nil, err
	}

	err = o.authz.CheckEvent(ctx, ident.Sub.ID, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	template := &models.BadgeTemplate{
		ID:           uuid.New(),
		EditionID:    input.EditionID,
		TicketTypeID: input.TicketTypeID,
		Name:         input.Name,
		DesignData:   input.DesignData,
	}

	return o.repo.Create(ctx, template)
}
