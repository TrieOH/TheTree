package service

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

// AwardStaffBadgesForEdition emits staff badges for every member of the
// edition's event. Called when an edition is published so staff added before
// publication still get their badge.
func (o *Operations) AwardStaffBadgesForEdition(ctx context.Context, editionID uuid.UUID) error {
	edition, err := o.Editions.GetByID(ctx, editionID)
	if err != nil {
		return err
	}

	members, err := o.Events.ListEventMembers(ctx, edition.EventID)
	if err != nil {
		return err
	}

	for _, member := range members {
		_, err = o.Emissions.Upsert(ctx, &models.BadgeEmission{
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
