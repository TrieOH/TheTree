package checkouts

import (
	"context"

	"lib/telemetry"

	"github.com/google/uuid"
)

// AttendeeCount returns the number of confirmed attendees of an edition —
// the public attendee-count read (GET /editions/{id}/attendees/count), the
// storefront's "N people already registered" number. Only paid
// registrations count (status confirmed); pending reservations and
// cancelled/expired rows do not. Unknown editions are NOT_FOUND (the
// endpoint exists, the edition does not).
func (o *Operations) AttendeeCount(ctx context.Context, editionID uuid.UUID) (int64, error) {
	ctx, span := telemetry.StartSpan(ctx, "CheckoutsService.AttendeeCount")
	defer span.End()

	_, err := o.editions.GetByID(ctx, editionID)
	if err != nil {
		return 0, err
	}
	return o.registrations.CountConfirmedByEdition(ctx, editionID)
}
