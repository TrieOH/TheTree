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

func (o *Operations) Close(ctx context.Context, formID uuid.UUID) (*models.Form, error) {
	ctx, span := telemetry.StartSpan(ctx, "FormService.Close")
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

	if form.Status != models.FormStatusOpen {
		return nil, fun.ErrBadRequest("cannot close a form not on open")
	}

	err = authz.Service.CheckForm(ctx, ident.Sub.ID, form.ID, models.FormMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.forms.Close(ctx, formID)
}
