package oauth_providers

import (
	"IdentityX/models"
	"context"
	"lib/crypto"
	"lib/telemetry"

	"go.uber.org/zap"
)

func (o *Operations) Update(ctx context.Context, payload models.UpdateOAuthProviderInput) (*models.ProjectOAuthProviders, error) {
	ctx, span := telemetry.StartSpan(ctx, "Update")
	defer span.End()

	row, err := o.providers.GetByID(ctx, payload.ID)
	if err != nil {
		return nil, err
	}
	err = o.requireProjectAdmin(ctx, row.ProjectID)
	if err != nil {
		return nil, err
	}

	if payload.ClientID != nil {
		row, err = o.providers.UpdateClientID(ctx, payload.ID, *payload.ClientID)
		if err != nil {
			return nil, err
		}
	}
	if payload.ClientSecret != nil {
		encryptedSecret, err := crypto.EncryptPrivateKey([]byte(*payload.ClientSecret))
		if err != nil {
			telemetry.Log().Error("encrypt oauth client secret", zap.Error(err))
			return nil, err
		}
		row, err = o.providers.UpdateClientSecret(ctx, payload.ID, encryptedSecret)
		if err != nil {
			return nil, err
		}
	}

	return row, nil
}
