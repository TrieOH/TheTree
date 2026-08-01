package forms

import (
	"context"
	idx "sdk/identityx"

	"Informd/models"
	"lib/telemetry"
)

func (o *Operations) ListArchivedForms(ctx context.Context) (forms []models.Form, err error) {
	ctx, span := telemetry.StartSpan(ctx, "FormService.ListArchivedForms")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	forms, err = o.forms.ListMineArchived(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	return forms, nil
}
