package forms

import (
	"context"
	idx "sdk/identityx"

	"Informd/models"
	"lib/telemetry"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (o *Operations) ReDraft(ctx context.Context, formID uuid.UUID) (*models.Form, error) {
	ctx, span := telemetry.StartSpan(ctx, "FormService.ReDraft")
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
		return nil, fun.ErrBadRequest("cannot redraft a form not on open")
	}

	count, err := o.forms.ResponsesCount(ctx, formID)
	if err != nil {
		return nil, err
	}

	if count != 0 {
		return nil, fun.ErrBadRequest("cannot redraft a form with responses")
	}

	err = o.authz.CheckForm(ctx, ident.Sub.ID, form.ID, models.FormMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.forms.ReDraft(ctx, formID)
}
