package commands

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/internal/authz"
	"univents/models"

	"github.com/google/uuid"
)

func (c *Commands) DeleteProgram(ctx context.Context, id uuid.UUID) (*models.Program, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsService.DeleteProgram")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	existing, err := c.programs.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	edition, err := c.editions.GetByID(ctx, existing.EditionID)
	if err != nil {
		return nil, err
	}

	err = authz.Service.CheckEvent(ctx, ident.Sub.ID, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return c.programs.Delete(ctx, id)
}
