package service

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

	form, err := o.Forms.GetByID(ctx, formID)
	if err != nil {
		return 0, err
	}

	err = o.Authz.CheckForm(ctx, ident.Sub.ID, form.ID, models.FormMemberRoleMember)
	if err != nil {
		return 0, err
	}

	return o.Forms.ResponsesCount(ctx, formID)
}
