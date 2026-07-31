package queries

import (
	"context"
	idx "sdk/identityx"

	"Informd/internal/authz"
	"Informd/models"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (q *Queries) GetResponseCount(ctx context.Context, formID uuid.UUID) (int, error) {
	ctx, span := telemetry.StartSpan(ctx, "FormService.GetResponseCount")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return 0, err
	}

	form, err := q.forms.GetByID(ctx, formID)
	if err != nil {
		return 0, err
	}

	err = authz.Service.CheckForm(ctx, ident.Sub.ID, form.ID, models.FormMemberRoleMember)
	if err != nil {
		return 0, err
	}

	return q.forms.ResponsesCount(ctx, formID)
}
