package service

import (
	"context"
	idx "sdk/identityx"

	"Informd/models"
)

func (o *Operations) ListArchivedForms(ctx context.Context) (forms []models.Form, err error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	forms, err = o.Forms.ListMineArchived(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	return forms, nil
}
