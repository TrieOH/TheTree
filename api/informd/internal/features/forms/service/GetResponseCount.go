package service

import (
	"context"

	"Informd/models"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (o *Operations) GetResponseCount(ctx context.Context, formID uuid.UUID) (int, error) {
	ctx, span := telemetry.StartSpan(ctx, "FormService.GetResponseCount")
	defer span.End()

	form, err := o.Forms.GetByID(ctx, formID)
	if err != nil {
		return 0, err
	}

	_, err = o.Authz.Form(ctx, form.ID, models.FormMemberRoleMember)
	if err != nil {
		return 0, err
	}

	return o.Forms.ResponsesCount(ctx, formID)
}
