package service

import (
	"context"

	"Informd/errs"
	"Informd/models"

	"github.com/google/uuid"
)

// TODO: kill this duplicated namespaced route — authz.Form already anchors via the form's namespace.
func (o *Operations) ReDraftForm(ctx context.Context, namespaceID, formID uuid.UUID) (*models.Form, error) {
	form, err := o.Forms.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}

	if form.Status != models.FormStatusOpen {
		return nil, errs.ErrRedraftNotOpen
	}

	count, err := o.Forms.ResponsesCount(ctx, formID)
	if err != nil {
		return nil, err
	}

	if count != 0 {
		return nil, errs.ErrRedraftWithResponses
	}

	_, err = o.Authz.Namespace(ctx, namespaceID, models.NamespaceMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.Forms.ReDraft(ctx, form.ID)
}
