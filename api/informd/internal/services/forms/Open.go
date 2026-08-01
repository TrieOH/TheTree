package forms

import (
	"context"
	idx "sdk/identityx"

	"Informd/internal/authz"
	"Informd/models"
	"lib/telemetry"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (o *Operations) Open(ctx context.Context, formID uuid.UUID) (*models.Form, error) {
	ctx, span := telemetry.StartSpan(ctx, "FormService.Open")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	var form *models.Form
	form, err = o.forms.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}

	if form.Status != models.FormStatusDraft {
		return nil, fun.ErrBadRequest("cannot open a form not on draft")
	}

	err = authz.Service.CheckForm(ctx, ident.Sub.ID, form.ID, models.FormMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.forms.Open(ctx, formID)
}
