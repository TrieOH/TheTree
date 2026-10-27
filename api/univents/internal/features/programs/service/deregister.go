package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"lib/database"
	"univents/errs"
	"univents/models"
)

// Deregister cancels the caller's sign-up for an occurrence (self-service).
// The spot is locked once attendance is marked: the guarded update flips
// only 'registered' rows, so an attended/no_show participation (or a racing
// state change) returns 0 rows → 409. Deregistering is a status flip, never
// a delete — the participation stays in the ledger as cancelled history.
func (o *Operations) Deregister(ctx context.Context, occurrenceID, userID uuid.UUID) (*models.ProgramParticipation, error) {
	occ, err := o.Occurrences.GetOccurrenceByID(ctx, occurrenceID)
	if err != nil {
		return nil, err
	}

	reg, err := o.Registrations.GetActiveByEditionAndAttendee(ctx, occ.EditionID, userID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, errs.ErrNotRegistered
		}
		return nil, err
	}

	part, err := o.Participations.GetActiveByOccurrenceAndRegistration(ctx, occurrenceID, reg.ID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, errs.ErrNotRegistered
		}
		return nil, err
	}

	updated, err := o.Participations.UpdateParticipationStatusIfRegistered(ctx, part.ID, models.ProgramParticipationStatusCancelled)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, errs.ErrRegistrationLocked
	}

	o.notifyStock(ctx, occ.EditionID, occurrenceID)
	return updated, nil
}
