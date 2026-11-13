package service

import (
	"context"
	"uuid"

	"Informd/models"
)

func (o *Operations) EditSelectConfig(ctx context.Context, formID uuid.UUID, payload models.FieldSelectConfig) (*models.FieldSelectConfig, error) {
	_, err := o.Authz.Form(ctx, formID, models.FormMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.Fields.UpdateSelectConfig(ctx, payload)
}
