package commands

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/internal/authz"
	"univents/models"
)

func (c *Commands) Create(ctx context.Context, payload models.AddSignatureInput) (*models.Signature, error) {
	ctx, span := telemetry.StartSpan(ctx, "SignaturesService.Create")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	edition, err := c.editions.GetByID(ctx, payload.EditionID)
	if err != nil {
		return nil, err
	}

	err = authz.Service.CheckEvent(ctx, ident.Sub.ID, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	sig, err := c.signatures.Create(ctx, &models.Signature{
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

	c.sendConfirmationEmail(ctx, sig)

	return sig, nil
}
