package service

import (
	"context"
	"univents/models"
)

func (o *Operations) CreateProgram(ctx context.Context, payload models.CreateProgramInput) (*models.Program, error) {
	edition, err := o.Editions.GetByID(ctx, payload.EditionID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
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

	return o.Programs.Create(ctx, program)
}
