package service

import (
	"context"

	"Informd/models"

	"github.com/google/uuid"
)

func (o *Operations) List(ctx context.Context, formID, stepID uuid.UUID) ([]models.Field, error) {
	_, err := o.Authz.Form(ctx, formID, models.FormMemberRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Fields.ListByStepID(ctx, stepID)
}
