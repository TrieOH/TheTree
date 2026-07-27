package commands

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (c *Commands) DeleteOccurrence(ctx context.Context, id uuid.UUID) (*models.ProgramOccurrence, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsService.DeleteOccurrence")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	existing, err := c.occurrences.GetOccurrenceByID(ctx, id)
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

	return c.occurrences.DeleteOccurrence(ctx, id)
}
