package queries

import (
	"Informd/models"
	"context"
	"lib/authz"
)

func (s *QueryService) ListForms(ctx context.Context) (forms []models.Form, err error) {
	ctx, span := s.tracer.Start(ctx, "FormService.ListForms")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	forms, err = s.forms.ListMine(ctx, sub.ID)
	if err != nil {
		return nil, err
	}

	return forms, nil
}
