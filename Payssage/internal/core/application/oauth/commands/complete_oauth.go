package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/plataform/telemetry"
	"TriePayments/internal/shared/errx"
	"context"
	"fmt"

	"go.uber.org/zap"
)

func (uc *CommandService) CompleteOAuth(ctx context.Context, provider, stateToken, code, redirectURI string) (string, error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.CompleteOAuth")
	defer span.End()

	oauthState, err := uc.oauthStates.Get(ctx, stateToken)
	if err != nil {
		return "", errx.Invalid("oauth_state").SetMessage("invalid or expired state")
	}

	if oauthState.Provider != provider {
		return "", errx.Invalid("oauth_state").SetMessage("provider mismatch")
	}

	p, err := uc.getProvider(provider)
	if err != nil {
		return "", err
	}

	credData, err := p.ExchangeCode(ctx, code, redirectURI)
	if err != nil {
		telemetry.Log().Info("Error exchanging codes", zap.Error(err))
		return "", errx.Internal("oauth").SetMessage(fmt.Sprintf("failed to exchange code: %s", err.Error()))
	}

	cred, err := uc.credentials.Create(ctx, domain.ProviderCredential{
		WorkspaceID: oauthState.WorkspaceID,
		Provider:    provider,
		DisplayName: credData.ProviderUserID,
		Credentials: credData,
	})
	if err != nil {
		return "", err
	}

	// if setup flow + marketplace, auto-create marketplace config
	if oauthState.Flow == domain.OAuthFlowSetup && oauthState.IsMarketplace {
		existing, err := uc.marketplace.Get(ctx, oauthState.WorkspaceID, cred.ID)
		if err != nil && !errx.IsKind(err, "not_found") {
			return "", err
		}
		if existing != nil {
			_, err = uc.marketplace.Update(ctx, domain.MarketplaceConfig{
				WorkspaceID:  oauthState.WorkspaceID,
				CredentialID: cred.ID,
				FeeBps:       oauthState.FeeBps,
			})
		} else {
			_, err = uc.marketplace.Create(ctx, domain.MarketplaceConfig{
				WorkspaceID:  oauthState.WorkspaceID,
				CredentialID: cred.ID,
				FeeBps:       oauthState.FeeBps,
			})
		}
		if err != nil {
			return "", err
		}
	}

	_ = uc.oauthStates.Delete(ctx, stateToken)

	switch oauthState.Flow {
	case domain.OAuthFlowSetup:
		return fmt.Sprintf("%s?provider=%s&status=success", oauthState.FinalRedirectURL, provider), nil
	case domain.OAuthFlowConnect:
		return fmt.Sprintf("%s?credential_id=%s", oauthState.FinalRedirectURL, cred.ID), nil
	default:
		return oauthState.FinalRedirectURL, nil
	}
}
