package badges

import (
	"context"
	"lib/telemetry"
	"time"
	"univents/models"

	"github.com/google/uuid"
)

// AwardStaffBadgesForUser emits a staff badge for the user on every published
// edition of the event that is current or future (ends_at >= now). Past
// editions are never awarded retroactively; draft editions are awarded when
// published (see AwardStaffBadgesForEdition).
func (o *Operations) AwardStaffBadgesForUser(ctx context.Context, eventID, userID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "BadgesService.AwardStaffBadgesForUser")
	defer span.End()

	editions, err := o.editions.ListPublic(ctx, eventID)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, edition := range editions {
		if edition.EndsAt.Before(now) {
			continue
		}
		_, err = o.emissions.Upsert(ctx, &models.BadgeEmission{
			EditionID: edition.ID,
			UserID:    userID,
			Origin:    models.BadgeEmissionOriginStaff,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
