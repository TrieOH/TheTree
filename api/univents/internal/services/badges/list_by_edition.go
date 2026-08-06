package badges

import (
	"context"
	"lib/telemetry"
	"univents/models"

	idx "sdk/identityx"

	"github.com/google/uuid"
)

// ListByEdition returns all emissions of an edition for the owner/admin to
// browse (all statuses; holder names are derived front-side).
func (o *Operations) ListByEdition(ctx context.Context, editionID uuid.UUID) ([]models.BadgeEditionEmission, error) {
	ctx, span := telemetry.StartSpan(ctx, "BadgesService.ListByEdition")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	edition, err := o.editions.GetByID(ctx, editionID)
	if err != nil {
		return nil, err
	}
	err = o.authz.CheckEvent(ctx, ident.Sub.ID, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	views, err := o.emissions.ListViewsByEdition(ctx, editionID)
	if err != nil {
		return nil, err
	}
	idx, err := o.loadTemplateIndex(ctx, views)
	if err != nil {
		return nil, err
	}

	items := make([]models.BadgeEditionEmission, 0, len(views))
	for _, v := range views {
		item := models.BadgeEditionEmission{
			ID:           v.ID,
			UserID:       v.UserID,
			Origin:       v.Origin,
			Status:       v.Status,
			StatusReason: v.StatusReason,
			TicketName:   v.TicketName,
			EmittedAt:    v.EmittedAt,
			UpdatedAt:    v.UpdatedAt,
		}
		if t := idx.match(v.EditionID, v.TicketTypeID); t != nil {
			id := t.ID
			name := t.Name
			item.TemplateID = &id
			item.TemplateName = &name
		}
		items = append(items, item)
	}
	return items, nil
}
