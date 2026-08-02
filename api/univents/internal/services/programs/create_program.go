package programs

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/models"
)

func (o *Operations) CreateProgram(ctx context.Context, payload models.CreateProgramInput) (*models.Program, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsService.CreateProgram")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	edition, err := o.editions.GetByID(ctx, payload.EditionID)
	if err != nil {
		return nil, err
	}

	err = o.authz.CheckEvent(ctx, ident.Sub.ID, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
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

	return o.programs.Create(ctx, program)
}
