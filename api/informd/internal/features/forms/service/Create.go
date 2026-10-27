package service

import (
	"context"
	idx "sdk/identityx"
	"time"

	"Informd/models"
)

func (o *Operations) Create(ctx context.Context, title string) (*models.Form, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	form, err := models.NewForm(nil, ident.Sub.ID, ident.Sub.ID, title)
	if err != nil {
		return nil, err
	}

	var created *models.Form
	err = o.Tx.WithinTx(ctx, func(ctx context.Context) error {
		created, err = o.Forms.Create(ctx, *form)
		if err != nil {
			return err
		}

		return o.Forms.AddMember(ctx, models.FormMember{
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
