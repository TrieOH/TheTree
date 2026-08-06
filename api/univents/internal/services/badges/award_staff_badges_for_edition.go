package badges

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

// AwardStaffBadgesForEdition emits staff badges for every member of the
// edition's event. Called when an edition is published so staff added before
// publication still get their badge.
func (o *Operations) AwardStaffBadgesForEdition(ctx context.Context, editionID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "BadgesService.AwardStaffBadgesForEdition")
	defer span.End()

	edition, err := o.editions.GetByID(ctx, editionID)
	if err != nil {
		return err
	}

	members, err := o.events.ListEventMembers(ctx, edition.EventID)
	if err != nil {
		return err
	}

	for _, member := range members {
		_, err = o.emissions.Upsert(ctx, &models.BadgeEmission{
			EditionID: edition.ID,
			UserID:    member.UserID,
			Origin:    models.BadgeEmissionOriginStaff,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
