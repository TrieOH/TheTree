package service

import (
	"context"

	"Informd/models"

	"github.com/google/uuid"
)

func (o *Operations) List(ctx context.Context, formID uuid.UUID) ([]models.Step, error) {
	_, err := o.Authz.Form(ctx, formID, models.FormMemberRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Steps.List(ctx, formID)
}
