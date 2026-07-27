package commands

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/models"

	"github.com/MintzyG/fun"
)

func (c *Commands) CreateProgram(ctx context.Context, payload models.CreateProgramInput) (*models.Program, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsService.CreateProgram")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	edition, err := c.editions.GetByID(ctx, payload.EditionID)
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

	program := &models.Program{
		EditionID:      payload.EditionID,
		Kind:           payload.Kind,
		Name:           payload.Name,
		Description:    payload.Description,
		MinAccessLevel: payload.MinAccessLevel,
		StaffOnly:      payload.StaffOnly,
		Price:          payload.Price,
	}

	return c.programs.Create(ctx, program)
}
