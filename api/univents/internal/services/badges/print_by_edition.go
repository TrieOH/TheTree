package badges

import (
	"context"
	"encoding/json"
	"lib/telemetry"
	"univents/models"

	idx "sdk/identityx"

	"github.com/google/uuid"
)

// PrintByEdition returns the print payload for the edition's active emissions,
// optionally filtered to specific emission ids. Owner/admin only.
func (o *Operations) PrintByEdition(ctx context.Context, editionID uuid.UUID, emissionIDs []uuid.UUID) ([]models.BadgePrintItem, error) {
	ctx, span := telemetry.StartSpan(ctx, "BadgesService.PrintByEdition")
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

	include := map[uuid.UUID]struct{}{}
	for _, id := range emissionIDs {
		include[id] = struct{}{}
	}

	items := make([]models.BadgePrintItem, 0, len(views))
	for _, v := range views {
		if v.Status != models.BadgeEmissionStatusActive {
			continue
		}
		if len(include) > 0 {
			if _, ok := include[v.ID]; !ok {
				continue
			}
		}
		item := models.BadgePrintItem{
			EmissionID:  v.ID,
			UserID:      v.UserID,
			Origin:      v.Origin,
			EventName:   v.EventName,
			EditionName: v.EditionName,
			TicketName:  v.TicketName,
			ActionURL:   profileURL(v.UserID),
		}
		if t := idx.match(v.EditionID, v.TicketTypeID, v.Origin); t != nil {
			id := t.ID
			name := t.Name
			item.TemplateID = &id
			item.TemplateName = &name
			item.DesignData = t.DesignData
		} else {
			item.DesignData = json.RawMessage(`{}`)
		}
		items = append(items, item)
	}
	return items, nil
}
