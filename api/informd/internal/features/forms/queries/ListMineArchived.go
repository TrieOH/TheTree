package queries

import (
	"context"

	"Informd/models"
	"lib/authz"
)

func (q *Queries) ListArchivedForms(ctx context.Context) (forms []models.Form, err error) {
	ctx, span := q.tracer.Start(ctx, "FormService.ListArchivedForms")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	forms, err = q.forms.ListMineArchived(ctx, sub.ID)
	if err != nil {
		return nil, err
	}

	return forms, nil
}
