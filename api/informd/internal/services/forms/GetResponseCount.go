package forms

import (
	"context"
	idx "sdk/identityx"

	"Informd/models"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (o *Operations) GetResponseCount(ctx context.Context, formID uuid.UUID) (int, error) {
	ctx, span := telemetry.StartSpan(ctx, "FormService.GetResponseCount")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return 0, err
	}

	form, err := o.forms.GetByID(ctx, formID)
	if err != nil {
		return 0, err
	}

	err = o.authz.CheckForm(ctx, ident.Sub.ID, form.ID, models.FormMemberRoleMember)
	if err != nil {
		return 0, err
	}

	return o.forms.ResponsesCount(ctx, formID)
}
