package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"lib/database"
	"univents/errs"
	"univents/models"
)

// CheckIn marks an attendee as present at a checkpoint occurrence
// (staff-only). Checkpoints are pass-through by design — attendees never
// register for them — so the gate is "can they pass", not "did they sign
// up": a confirmed edition ticket plus the checkpoint's access rules
// (staff_only, min_access_level), resolved from the attendee's badge QR.
// The first scan upserts an attended participation anchored to the
// attendee's edition-ticket registration; re-scanning is idempotent (an
// existing live row is flipped to attended, never duplicated — the
// partial unique index uniq_program_participations_active_per_occurrence_attendee
// is the backstop). Activities keep the participation-driven
// mark-attended flow and are rejected here (400).
func (o *Operations) CheckIn(ctx context.Context, occurrenceID, attendeeID, actorID uuid.UUID) (*models.ProgramParticipation, error) {
	occ, err := o.Occurrences.GetOccurrenceByID(ctx, occurrenceID)
	if err != nil {
		return nil, err
	}

	edition, err := o.Editions.GetByID(ctx, occ.EditionID)
	if err != nil {
		return nil, err
	}
	err = o.Authz.EventFor(ctx, actorID, edition.EventID, models.EventMemberRoleStaff)
	if err != nil {
		return nil, err
	}

	program, err := o.Programs.GetByID(ctx, occ.ProgramID)
	if err != nil {
		return nil, err
	}
	if program.Kind != models.ProgramKindCheckpoint {
		return nil, errs.ErrCheckInNotCheckpoint
	}

	reg, err := o.Registrations.GetActiveByEditionAndAttendee(ctx, occ.EditionID, attendeeID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, errs.ErrTicketRequiredForCheckpoint
		}
		return nil, err
	}
	if reg.Status != models.RegistrationStatusConfirmed {
		return nil, errs.ErrTicketRequiredForCheckpoint
	}
	err = o.checkProgramAccess(ctx, occ, program, reg, attendeeID)
	if err != nil {
		return nil, err
	}

	return o.Participations.UpsertAttended(ctx, occ.EditionID, occurrenceID, reg.ID)
}
