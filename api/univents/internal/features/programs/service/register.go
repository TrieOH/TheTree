package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"lib/database"
	"univents/errs"
	"univents/models"
)

// Register signs the caller up for an occurrence (self-service). Flow:
// resolve + permission-check (confirmed ticket, staff_only, access level),
// then inside a tx: row-lock the occurrence (serializes concurrent
// sign-ups), 409 on double registration, 409 when full, insert. The partial
// unique index
// (uniq_program_participations_active_per_occurrence_attendee) is the
// atomic backstop for two racing inserts. On success the occurrence's stock
// delta is published so the store's realtime surfaces stay fresh.
func (o *Operations) Register(ctx context.Context, occurrenceID, userID uuid.UUID) (*models.ProgramParticipation, error) {
	occ, err := o.Occurrences.GetOccurrenceByID(ctx, occurrenceID)
	if err != nil {
		return nil, err
	}

	reg, _, err := o.checkEligible(ctx, occ, userID)
	if err != nil {
		return nil, err
	}

	var part *models.ProgramParticipation
	err = o.Tx.WithinTx(ctx, func(ctx context.Context) error {
		// Lock the occurrence row: two concurrent sign-ups serialize here,
		// so the count below cannot oversell the last slot.
		_, err := o.Occurrences.GetOccurrenceByIDForUpdate(ctx, occurrenceID)
		if err != nil {
			return err
		}

		existing, err := o.Participations.GetActiveByOccurrenceAndRegistration(ctx, occurrenceID, reg.ID)
		if existing != nil {
			return errs.ErrAlreadyRegistered
		}
		if err != nil && !errors.Is(err, database.ErrNotFound) {
			return err
		}

		count, err := o.Participations.CountActiveByOccurrence(ctx, occurrenceID)
		if err != nil {
			return err
		}
		if occ.MaxCapacity != nil && count >= int64(*occ.MaxCapacity) {
			return errs.ErrOccurrenceFull
		}

		part, err = o.Participations.CreateParticipation(ctx, &models.ProgramParticipation{
			EditionID:      occ.EditionID,
			OccurrenceID:   occurrenceID,
			RegistrationID: reg.ID,
			Status:         models.ProgramParticipationStatusRegistered,
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	o.notifyStock(ctx, occ.EditionID, occurrenceID)
	return part, nil
}
