package commands

import (
	"TriePayments/internal/core/domain"
	"context"
	"fmt"

	"TriePayments/internal/shared/errx"
)

func (uc *CommandService) CompleteOAuth(ctx context.Context, provider, stateToken, code string) (string, error) {
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

	credData, err := p.ExchangeCode(ctx, code)
	if err != nil {
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

	if err := uc.oauthStates.Delete(ctx, stateToken); err != nil {
		// non-fatal, state will expire naturally
		_ = err
	}

	finalURL := fmt.Sprintf("%s?credential_id=%s", oauthState.FinalRedirectURL, cred.ID)
	return finalURL, nil
}
