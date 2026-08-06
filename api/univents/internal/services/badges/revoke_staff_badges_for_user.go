package badges

import (
	"context"
	"lib/telemetry"
	"time"
	"univents/models"

	"github.com/google/uuid"
)

// RevokeStaffBadgesForUser revokes the user's staff badges on current and
// future editions of the event. Past-edition badges are kept as keepsakes.
func (o *Operations) RevokeStaffBadgesForUser(ctx context.Context, eventID, userID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "BadgesService.RevokeStaffBadgesForUser")
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
		err = o.emissions.Revoke(ctx, edition.ID, userID, models.BadgeEmissionOriginStaff, "staff removed")
		if err != nil {
			return err
		}
	}
	return nil
}
