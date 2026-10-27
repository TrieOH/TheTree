package service

import (
	"context"

	"Informd/errs"
	"Informd/models"

	"github.com/google/uuid"
)

// TODO: kill this duplicated namespaced route — authz.Form already anchors via the form's namespace.
func (o *Operations) ArchiveForm(ctx context.Context, namespaceID, formID uuid.UUID) (*models.Form, error) {
	form, err := o.Forms.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}

	if form.Status != models.FormStatusClosed {
		return nil, errs.ErrArchiveNotClosed
	}

	_, err = o.Authz.Namespace(ctx, namespaceID, models.NamespaceMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.Forms.Close(ctx, form.ID)
}
