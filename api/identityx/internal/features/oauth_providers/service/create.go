package service

import (
	"IdentityX/models"
	"context"
	"lib/telemetry"
)

func (o *Operations) Create(ctx context.Context, payload models.CreateOAuthProviderInput) (*models.ProjectOAuthProviders, error) {
	err := o.Authz.CheckProject(ctx, payload.ProjectID, models.ProjectRoleAdmin)
	if err != nil {
		return nil, err
	}

	encryptedSecret, err := o.Vault.EncryptPrivateKey([]byte(payload.ClientSecret))
	if err != nil {
		telemetry.Log().Error("encrypt oauth client secret", telemetry.Err(err))
		return nil, err
	}

	created, err := o.Providers.Create(ctx, models.ProjectOAuthProviders{
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
