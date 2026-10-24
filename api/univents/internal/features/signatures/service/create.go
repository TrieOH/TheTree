package service

import (
	"context"
	"lib/telemetry"
	"univents/models"
)

func (o *Operations) Create(ctx context.Context, payload models.AddSignatureInput) (*models.Signature, error) {
	ctx, span := telemetry.StartSpan(ctx, "SignaturesService.Create")
	defer span.End()

	edition, err := o.Editions.GetByID(ctx, payload.EditionID)
	if err != nil {
		return nil, err
	}

	ident, err := o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	sig, err := o.Signatures.Create(ctx, &models.Signature{
		EditionID:       payload.EditionID,
		CreatedBy:       ident.Sub.ID,
		SignatoryName:   payload.SignatoryName,
		SignatoryTitle:  payload.SignatoryTitle,
		SignatoryEmail:  payload.SignatoryEmail,
		SignatoryUserID: payload.SignatoryUserID,
		ImageURL:        payload.ImageURL,
	})
	if err != nil {
		return nil, err
	}

	o.sendConfirmationEmail(ctx, sig)

	return sig, nil
}
