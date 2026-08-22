package programs

import (
	"context"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"

	"univents/models"
)

// MarkAttended flips a participation to attended (staff-only). Any event
// staff (owner/admin/staff) may mark any activity in the edition — the
// domain's per-program staff assignment is deferred. The guarded update
// flips any non-cancelled row: re-marking an attended participation is
// idempotent, a no_show→attended correction is allowed, and a cancelled
// participation (refunded purchase, or the attendee dropped out) is 409.
func (o *Operations) MarkAttended(ctx context.Context, participationID, actorID uuid.UUID) (*models.ProgramParticipation, error) {
	part, err := o.participations.GetParticipationByID(ctx, participationID)
	if err != nil {
		return nil, err
	}

	edition, err := o.editions.GetByID(ctx, part.EditionID)
	if err != nil {
		return nil, err
	}
	err = o.authz.CheckEvent(ctx, actorID, edition.EventID, models.EventMemberRoleStaff)
	if err != nil {
		return nil, err
	}

	updated, err := o.participations.UpdateParticipationStatusIfNotCancelled(ctx, participationID, models.ProgramParticipationStatusAttended)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, fun.ErrConflict("cannot mark attendance: participation was cancelled")
	}
	return updated, nil
}
