package queries

import (
	"context"
	idx "sdk/identityx"

	"Informd/models"
	"lib/telemetry"
)

func (q *Queries) ListForms(ctx context.Context) (forms []models.Form, err error) {
	ctx, span := telemetry.StartSpan(ctx, "FormService.ListForms")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	forms, err = q.forms.ListMine(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	return forms, nil
}
