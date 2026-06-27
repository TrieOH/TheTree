package commands

import (
	"context"
	idx "sdk/identityx"
	"time"

	"Informd/models"
)

func (s *Commands) Create(ctx context.Context, title string) (*models.Form, error) {
	ctx, span := s.tracer.Start(ctx, "FormService.Create")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	form, err := models.NewForm(nil, ident.Sub.ID, ident.Sub.ID, title)
	if err != nil {
		return nil, err
	}

	var created *models.Form
	if err = s.tx.WithinTx(ctx, func(ctx context.Context) error {
		created, err = s.forms.Create(ctx, *form)
		if err != nil {
			return err
		}

		return s.forms.AddMember(ctx, models.FormMember{
			UserID:  ident.Sub.ID,
			FormID:  created.ID,
			Role:    models.FormMemberRoleOwner,
			AddedAt: time.Now(),
			AddedBy: ident.Sub.ID,
		})
	}); err != nil {
		return nil, err
	}

	return created, nil
}
