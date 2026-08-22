package programs

import (
	"context"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"

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
	occ, err := o.occurrences.GetOccurrenceByID(ctx, occurrenceID)
	if err != nil {
		return nil, err
	}

	reg, _, err := o.checkEligible(ctx, occ, userID)
	if err != nil {
		return nil, err
	}

	var part *models.ProgramParticipation
	err = o.tx.WithinTx(ctx, func(ctx context.Context) error {
		// Lock the occurrence row: two concurrent sign-ups serialize here,
		// so the count below cannot oversell the last slot.
		_, err := o.occurrences.GetOccurrenceByIDForUpdate(ctx, occurrenceID)
		if err != nil {
			return err
		}

		existing, err := o.participations.GetActiveByOccurrenceAndRegistration(ctx, occurrenceID, reg.ID)
		if existing != nil {
			return fun.ErrConflict("you are already registered for this activity")
		}
		if err != nil && !fun.Is(err, fun.CodeNotFound) {
			return err
		}

		count, err := o.participations.CountActiveByOccurrence(ctx, occurrenceID)
		if err != nil {
			return err
		}
		if occ.MaxCapacity != nil && count >= int64(*occ.MaxCapacity) {
			return fun.ErrConflict("this activity occurrence is full")
		}

		part, err = o.participations.CreateParticipation(ctx, &models.ProgramParticipation{
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
