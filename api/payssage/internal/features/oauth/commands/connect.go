package commands

import (
	"context"
	"payssage/internal/providers"
	"payssage/models"
	idx "sdk/identityx"
	"time"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (c *Commands) Connect(ctx context.Context, payload models.ConnectInput) (string, error) {
	ctx, span := c.tracer.Start(ctx, "ConnectCollector")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return "", err
	}

	var walletID *uuid.UUID
	if payload.Flow == models.OAuthFlowSeller {
		if ident.Cred.Type != idx.ApiKeyCredentialType {
			return "", fun.ErrForbidden("only api keys may connect sellers")
		}
		var wallet *models.Wallet
		wallet, err = c.wallets.GetByID(ctx, *payload.WalletID)
		if err != nil {
			return "", err
		}
		var org *models.Organization
		if wallet.OrganizationID != nil {
			org, err = c.orgs.GetByID(ctx, *wallet.OrganizationID)
			if err != nil {
				return "", err
			}
		}
		if org != nil {
			err = c.checkRole(ctx, org, ident.Sub.ID, models.OrganizationRoleAdmin)
		} else {
			err = c.checkWalletAccess(wallet, ident.Sub.ID)
		}
		if err != nil {
			return "", err
		}
		walletID = &wallet.ID
	}

	provider, err := providers.FromString(payload.Provider)
	if err != nil {
		return "", err
	}

	stateToken, err := generateState()
	if err != nil {
		return "", err
	}

	_, err = c.oauth.Create(ctx, models.OAuthState{
		State:            stateToken,
		WalletID:         walletID,
		Provider:         provider.String(),
		Flow:             payload.Flow,
		FinalRedirectUrl: payload.FinalRedirectURL,
		ExpiresAt:        time.Now().Add(5 * time.Minute),
	})
	if err != nil {
		return "", err
	}

	return providers.PayssageProviders.OAuth[provider].BuildAuthURL(stateToken, payload.ProviderRedirectURL), nil
}
