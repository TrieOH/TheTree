package programs

import (
	"context"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"

	"univents/models"
)

// Deregister cancels the caller's sign-up for an occurrence (self-service).
// The spot is locked once attendance is marked: the guarded update flips
// only 'registered' rows, so an attended/no_show participation (or a racing
// state change) returns 0 rows → 409. Deregistering is a status flip, never
// a delete — the participation stays in the ledger as cancelled history.
func (o *Operations) Deregister(ctx context.Context, occurrenceID, userID uuid.UUID) (*models.ProgramParticipation, error) {
	occ, err := o.occurrences.GetOccurrenceByID(ctx, occurrenceID)
	if err != nil {
		return nil, err
	}

	reg, err := o.registrations.GetActiveByEditionAndAttendee(ctx, occ.EditionID, userID)
	if err != nil {
		if fun.Is(err, fun.CodeNotFound) {
			return nil, fun.ErrNotFound("you are not registered for this activity")
		}
		return nil, err
	}

	part, err := o.participations.GetActiveByOccurrenceAndRegistration(ctx, occurrenceID, reg.ID)
	if err != nil {
		if fun.Is(err, fun.CodeNotFound) {
			return nil, fun.ErrNotFound("you are not registered for this activity")
		}
		return nil, err
	}

	updated, err := o.participations.UpdateParticipationStatusIfRegistered(ctx, part.ID, models.ProgramParticipationStatusCancelled)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, fun.ErrConflict("registration is locked: attendance has been marked")
	}

	o.notifyStock(ctx, occ.EditionID, occurrenceID)
	return updated, nil
}
