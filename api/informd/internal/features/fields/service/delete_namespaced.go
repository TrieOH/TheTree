package service

import (
	"context"

	"Informd/models"

	"github.com/google/uuid"
)

// TODO: kill this duplicated namespaced route — authz.Form already anchors via the form's namespace.
func (o *Operations) DeleteNamespaced(ctx context.Context, _, formID, fieldID uuid.UUID) error {
	_, err := o.Authz.Form(ctx, formID, models.FormMemberRoleAdmin)
	if err != nil {
		return err
	}

	return o.Fields.Delete(ctx, fieldID)
}
