package programs

import (
	"context"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"

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
	occ, err := o.occurrences.GetOccurrenceByID(ctx, occurrenceID)
	if err != nil {
		return nil, err
	}

	edition, err := o.editions.GetByID(ctx, occ.EditionID)
	if err != nil {
		return nil, err
	}
	err = o.authz.CheckEvent(ctx, actorID, edition.EventID, models.EventMemberRoleStaff)
	if err != nil {
		return nil, err
	}

	program, err := o.programs.GetByID(ctx, occ.ProgramID)
	if err != nil {
		return nil, err
	}
	if program.Kind != models.ProgramKindCheckpoint {
		return nil, fun.Err("check-in is only for checkpoint occurrences").BadRequest()
	}

	reg, err := o.registrations.GetActiveByEditionAndAttendee(ctx, occ.EditionID, attendeeID)
	if err != nil {
		if fun.Is(err, fun.CodeNotFound) {
			return nil, fun.ErrForbidden("a confirmed ticket is required to pass this checkpoint")
		}
		return nil, err
	}
	if reg.Status != models.RegistrationStatusConfirmed {
		return nil, fun.ErrForbidden("a confirmed ticket is required to pass this checkpoint")
	}
	err = o.checkProgramAccess(ctx, occ, program, reg, attendeeID)
	if err != nil {
		return nil, err
	}

	return o.participations.UpsertAttended(ctx, occ.EditionID, occurrenceID, reg.ID)
}
