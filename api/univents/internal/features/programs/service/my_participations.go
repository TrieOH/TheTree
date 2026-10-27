package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"lib/database"
	"univents/models"
)

// MyParticipations returns the caller's live activity sign-ups in the
// edition ("my activities"), joined with program + occurrence. Returns
// (nil, nil) when the caller holds no ticket — no sign-ups is a normal
// state, like the my-ticket read. Cancelled rows are history and never
// shown: the append-only ledger's current-state projection.
func (o *Operations) MyParticipations(ctx context.Context, editionID, userID uuid.UUID) ([]models.MyParticipation, error) {
	_, err := o.Editions.GetByID(ctx, editionID)
	if err != nil {
		return nil, err
	}

	reg, err := o.Registrations.GetActiveByEditionAndAttendee(ctx, editionID, userID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			// no ticket → no sign-ups is a normal state, not an error
			return nil, nil
		}
		return nil, err
	}

	return o.Participations.ListActiveByEditionAndRegistration(ctx, editionID, reg.ID)
}
