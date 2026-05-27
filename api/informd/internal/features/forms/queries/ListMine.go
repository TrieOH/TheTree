package queries

import (
	models2 "IdentityX/models"
	"Informd/models"
	"context"
)

func (s *QueryService) ListForms(ctx context.Context) (forms []models.Form, err error) {
	ctx, span := s.tracer.Start(ctx, "FormService.ListForms")
	defer span.End()

	var sub *models2.UserSubject
	sub, err = models2.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	forms, err = s.forms.ListMine(ctx, sub.ID)
	if err != nil {
		return nil, err
	}

	return forms, nil
}
