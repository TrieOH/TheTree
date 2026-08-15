package oauth_providers

import (
	"IdentityX/models"
	"context"
	"lib/crypto"
	"lib/telemetry"

	"go.uber.org/zap"
)

func (o *Operations) Create(ctx context.Context, payload models.CreateOAuthProviderInput) (*models.ProjectOAuthProviders, error) {
	ctx, span := telemetry.StartSpan(ctx, "Create")
	defer span.End()

	err := o.authz.CheckProject(ctx, payload.ProjectID, models.ProjectRoleAdmin)
	if err != nil {
		return nil, err
	}

	encryptedSecret, err := crypto.EncryptPrivateKey([]byte(payload.ClientSecret))
	if err != nil {
		telemetry.Log().Error("encrypt oauth client secret", zap.Error(err))
		return nil, err
	}

	created, err := o.providers.Create(ctx, models.ProjectOAuthProviders{
		ProjectID:             payload.ProjectID,
		Provider:              payload.Provider,
		ClientID:              payload.ClientID,
		EncryptedClientSecret: encryptedSecret,
		CallbackURL:           payload.CallbackURL,
		Enabled:               true,
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}
