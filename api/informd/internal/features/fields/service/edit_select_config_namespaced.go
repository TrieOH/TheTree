package service

import (
	"context"

	"Informd/models"

	"github.com/google/uuid"
)

// TODO: kill this duplicated namespaced route — authz.Form already anchors via the form's namespace.
func (o *Operations) EditSelectConfigNamespaced(ctx context.Context, formID, _ uuid.UUID, payload models.FieldSelectConfig) (*models.FieldSelectConfig, error) {
	_, err := o.Authz.Form(ctx, formID, models.FormMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.Fields.UpdateSelectConfig(ctx, payload)
}
