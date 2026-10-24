package service

import (
	"context"
	"univents/models"
)

func (o *Operations) PatchProgram(ctx context.Context, payload models.PatchProgramInput) (*models.Program, error) {
	existing, err := o.Programs.GetByID(ctx, payload.ProgramID)
	if err != nil {
		return nil, err
	}

	edition, err := o.Editions.GetByID(ctx, existing.EditionID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
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
