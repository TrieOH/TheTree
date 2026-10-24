package service

import (
	"context"

	"Informd/models"

	"github.com/google/uuid"
)

func (o *Operations) EditSelectConfig(ctx context.Context, formID uuid.UUID, payload models.FieldSelectConfig) (*models.FieldSelectConfig, error) {
	_, err := o.Authz.Form(ctx, formID, models.FormMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.Fields.UpdateSelectConfig(ctx, payload)
}
