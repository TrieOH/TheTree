package service

import (
	"context"
	"univents/models"
)

func (o *Operations) CreateTemplate(ctx context.Context, input models.CreateCertificationTemplateInput) (*models.CertificationTemplate, error) {
	edition, err := o.Editions.GetByID(ctx, input.EditionID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.Certs.CreateTemplate(ctx, input)
}
