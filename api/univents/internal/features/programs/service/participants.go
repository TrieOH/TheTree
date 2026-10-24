package service

import (
	"context"

	"github.com/google/uuid"

	"univents/models"
)

// Participants returns the attendance surface of an occurrence: every live
// participation with the attendee's identity (staff-only, owner/admin/staff
// — the same role the mark-attended write requires).
func (o *Operations) Participants(ctx context.Context, occurrenceID, actorID uuid.UUID) ([]models.ProgramParticipant, error) {
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

	return o.Participations.ListByOccurrence(ctx, occurrenceID)
}
