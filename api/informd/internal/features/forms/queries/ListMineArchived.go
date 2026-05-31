package queries

import (
	"Informd/models"
	"context"
	"lib/authz"
)

func (s *QueryService) ListArchivedForms(ctx context.Context) (forms []models.Form, err error) {
	ctx, span := s.tracer.Start(ctx, "FormService.ListArchivedForms")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	forms, err = s.forms.ListMineArchived(ctx, sub.ID)
	if err != nil {
		return nil, err
	}

	return forms, nil
}
