package checkouts

import (
	"context"

	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

// MyTicket returns the caller's active ticket for the edition — their own
// registration where attendee_user_id = caller and status is pending or
// confirmed — with its ticket type. Returns (nil, nil) when the caller
// holds none (cancelled/expired registrations do not count). Read-only;
// never transitions state. An unknown edition is NOT_FOUND (the endpoint
// exists, the edition does not).
func (o *Operations) MyTicket(ctx context.Context, editionID, userID uuid.UUID) (*models.MyTicket, error) {
	ctx, span := telemetry.StartSpan(ctx, "CheckoutsService.MyTicket")
	defer span.End()

	_, err := o.editions.GetByID(ctx, editionID)
	if err != nil {
		return nil, err
	}

	reg, err := o.activeRegistrationFor(ctx, editionID, userID)
	if err != nil {
		return nil, err
	}
	if reg == nil {
		//nolint:nilnil // no ticket is a normal state (front falls back to the buy flow)
		return nil, nil
	}

	ticketType, err := o.ticketTypes.GetByID(ctx, reg.TicketTypeID)
	if err != nil {
		return nil, err
	}
	return &models.MyTicket{
		RegistrationID: reg.ID,
		TicketType:     *ticketType,
		Status:         reg.Status,
	}, nil
}
