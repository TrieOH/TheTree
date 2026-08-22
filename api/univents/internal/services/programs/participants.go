package programs

import (
	"context"

	"github.com/google/uuid"

	"univents/models"
)

// Participants returns the attendance surface of an occurrence: every live
// participation with the attendee's identity (staff-only, owner/admin/staff
// — the same role the mark-attended write requires).
func (o *Operations) Participants(ctx context.Context, occurrenceID, actorID uuid.UUID) ([]models.ProgramParticipant, error) {
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

	return o.participations.ListByOccurrence(ctx, occurrenceID)
}
