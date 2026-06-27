package queries

import (
	"context"
	idx "sdk/identityx"

	"Informd/models"
)

func (q *Queries) ListArchivedForms(ctx context.Context) (forms []models.Form, err error) {
	ctx, span := q.tracer.Start(ctx, "FormService.ListArchivedForms")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	forms, err = q.forms.ListMineArchived(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	return forms, nil
}
