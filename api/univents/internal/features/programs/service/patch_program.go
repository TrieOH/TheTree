package service

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/models"
)

func (o *Operations) PatchProgram(ctx context.Context, payload models.PatchProgramInput) (*models.Program, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProgramsService.PatchProgram")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	existing, err := o.Programs.GetByID(ctx, payload.ProgramID)
	if err != nil {
		return nil, err
	}

	edition, err := o.Editions.GetByID(ctx, existing.EditionID)
	if err != nil {
		return nil, err
	}

	err = o.Authz.CheckEvent(ctx, ident.Sub.ID, edition.EventID, models.EventMemberRoleAdmin)
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

	return o.Programs.Patch(ctx, payload.ProgramID, program)
}
