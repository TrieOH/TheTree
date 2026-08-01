package commands

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/internal/authz"
	"univents/models"
)

func (c *Commands) PatchProgram(ctx context.Context, payload models.PatchProgramInput) (*models.Program, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsService.PatchProgram")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	existing, err := c.programs.GetByID(ctx, payload.ProgramID)
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

	program := &models.Program{
		Kind:           payload.Kind,
		Name:           payload.Name,
		Description:    payload.Description,
		MinAccessLevel: payload.MinAccessLevel,
		StaffOnly:      payload.StaffOnly,
		Price:          payload.Price,
		BannerURL:      payload.BannerURL,
	}

	return c.programs.Patch(ctx, payload.ProgramID, program)
}
