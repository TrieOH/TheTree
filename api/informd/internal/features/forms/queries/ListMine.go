package queries

import (
	"Informd/models"
	"context"
	"lib/authz"
)

func (q *Queries) ListForms(ctx context.Context) (forms []models.Form, err error) {
	ctx, span := q.tracer.Start(ctx, "FormService.ListForms")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	forms, err = q.forms.ListMine(ctx, sub.ID)
	if err != nil {
		return nil, err
	}

	return forms, nil
}
