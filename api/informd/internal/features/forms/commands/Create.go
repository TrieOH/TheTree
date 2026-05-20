package commands

import (
	"Informd/models"
	"context"
	"lib/authz"
	"time"
)

func (s *CommandService) Create(ctx context.Context, title string) (*models.Form, error) {
	ctx, span := s.tracer.Start(ctx, "FormService.Create")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	form, err := models.NewForm(nil, sub.ID, title)
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
			UserID:  sub.ID,
			FormID:  created.ID,
			Role:    models.FormMemberRoleOwner,
			AddedAt: time.Now(),
			AddedBy: sub.ID,
		})
	}); err != nil {
		return nil, err
	}

	return created, nil
}
