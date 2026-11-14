package service

import (
	"context"
	"uuid"

	"Informd/models"
)

func (o *Operations) Delete(ctx context.Context, formID, fieldID uuid.UUID) error {
	_, err := o.Authz.Form(ctx, formID, models.FormMemberRoleAdmin)
	if err != nil {
		return err
	}

	return o.Fields.Delete(ctx, fieldID)
}
