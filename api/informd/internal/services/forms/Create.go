package forms

import (
	"context"
	idx "sdk/identityx"
	"time"

	"Informd/models"
	"lib/database"
	"lib/telemetry"
)

func (o *Operations) Create(ctx context.Context, title string) (*models.Form, error) {
	ctx, span := telemetry.StartSpan(ctx, "FormService.Create")
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
	err = database.RunTx(ctx, func(ctx context.Context) error {
		created, err = o.forms.Create(ctx, *form)
		if err != nil {
			return err
		}

		return o.forms.AddMember(ctx, models.FormMember{
			UserID:  ident.Sub.ID,
			FormID:  created.ID,
			Role:    models.FormMemberRoleOwner,
			AddedAt: time.Now(),
			AddedBy: ident.Sub.ID,
		})
	})
	if err != nil {
		return nil, err
	}

	return created, nil
}
