package service

import (
	"context"

	"Informd/models"

	"github.com/google/uuid"
)

// TODO: kill this duplicated namespaced route — authz.Form already anchors via the form's namespace.
func (o *Operations) GetSelectConfigNamespaced(ctx context.Context, formID, _, fieldID uuid.UUID) (*models.FieldSelectConfig, error) {
	_, err := o.Authz.Form(ctx, formID, models.FormMemberRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Fields.GetSelectConfig(ctx, fieldID)
}
