package commands

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/models"

	"github.com/MintzyG/fun"
)

func (c *Commands) PatchOccurrence(ctx context.Context, payload models.PatchProgramOccurrenceInput) (*models.ProgramOccurrence, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsService.PatchOccurrence")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	existing, err := c.occurrences.GetOccurrenceByID(ctx, payload.OccurrenceID)
	if err != nil {
		return nil, err
	}

	edition, err := c.editions.GetByID(ctx, existing.EditionID)
	if err != nil {
		return nil, err
	}

	member, err := c.events.GetMember(ctx, edition.EventID, ident.Sub.ID)
	if fun.Is(err, fun.CodeNotFound) {
		return nil, fun.ErrForbidden("insufficient permissions")
	}
	if err != nil {
		return nil, err
	}
	if !member.Role.Minimum(models.EventMemberRoleAdmin) {
		return nil, fun.ErrForbidden("insufficient permissions")
	}

	occurrence := &models.ProgramOccurrence{
		StartsAt:    payload.StartsAt,
		EndsAt:      payload.EndsAt,
		MaxCapacity: payload.MaxCapacity,
	}

	return c.occurrences.PatchOccurrence(ctx, payload.OccurrenceID, occurrence)
}
