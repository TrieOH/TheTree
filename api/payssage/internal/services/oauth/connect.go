package oauth

import (
	"context"
	"lib/telemetry"
	"payssage/internal/providers"
	"payssage/models"
	idx "sdk/identityx"
	"time"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (o *Operations) Connect(ctx context.Context, payload models.ConnectInput) (string, error) {
	ctx, span := telemetry.StartSpan(ctx, "Connect")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return "", err
	}

	var walletID *uuid.UUID
	var organizationID *uuid.UUID

	switch payload.Flow {
	case models.OAuthFlowSeller:
		var wallet *models.Wallet
		wallet, err = o.wallets.GetByID(ctx, *payload.WalletID)
		if err != nil {
			return "", err
		}
		var org *models.Organization
		if wallet.OrganizationID != nil {
			org, err = o.orgs.GetByID(ctx, *wallet.OrganizationID)
			if err != nil {
				return "", err
			}
		}
		if org != nil {
			err = o.authz.CheckOrg(ctx, ident.Sub.ID, org.ID, models.OrganizationRoleAdmin)
		} else {
			err = o.authz.CheckWalletAccess(ctx, ident.Sub.ID, wallet.ID, models.OrganizationRoleMember)
		}
		if err != nil {
			return "", err
		}
		walletID = &wallet.ID

	case models.OAuthFlowCollector:
		if payload.OrganizationID != nil {
			var org *models.Organization
			org, err = o.orgs.GetByID(ctx, *payload.OrganizationID)
			if err != nil {
				return "", err
			}
			err = o.authz.CheckOrg(ctx, ident.Sub.ID, org.ID, models.OrganizationRoleAdmin)
			if err != nil {
				return "", err
			}
			organizationID = payload.OrganizationID
		}

	default:
		return "", fun.ErrBadRequest("invalid flow")
	}

	provider, err := providers.FromString(payload.Provider)
	if err != nil {
		return "", err
	}

	stateToken, err := generateState()
	if err != nil {
		return "", err
	}

	_, err = o.oauth.Create(ctx, models.OAuthState{
		State:            stateToken,
		WalletID:         walletID,
		OrganizationID:   organizationID,
		OwnerID:          ident.Sub.ID,
		Provider:         provider.String(),
		Flow:             payload.Flow,
		FinalRedirectURL: payload.FinalRedirectURL,
		ExpiresAt:        time.Now().Add(5 * time.Minute),
	})
	if err != nil {
		return "", err
	}

	// The provider redirect URI is Payssage's own callback (config-owned,
	// D7) — the provider builds the auth URL with it internally.
	return providers.PayssageProviders.OAuth[provider].BuildAuthURL(stateToken), nil
}
