package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/authz"
	"TriePayments/internal/shared/errx"
	"context"
	"time"
)

type ConnectSellerRequest struct {
	WorkspaceName       string
	Provider            string
	ProviderRedirectURL string
	FinalRedirectURL    string
}

func (uc *CommandService) ConnectSeller(ctx context.Context, req ConnectSellerRequest) (string, string, error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.ConnectSeller")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return "", "", err
	}

	workspace, err := uc.workspaces.GetByName(ctx, req.WorkspaceName, sub.ID)
	if err != nil {
		return "", "", err
	}

	_, err = uc.marketplace.GetByProvider(ctx, workspace.ID, req.Provider)
	if err != nil {
		return "", "", err
	}

	stateToken, err := generateState()
	if err != nil {
		return "", "", errx.Internal("oauth_state").SetCause(err)
	}

	_, err = uc.oauthStates.Create(ctx, domain.OAuthState{
		State:            stateToken,
		WorkspaceID:      workspace.ID,
		Provider:         req.Provider,
		Flow:             domain.OAuthFlowConnect,
		IsMarketplace:    false,
		FeeBps:           0,
		FinalRedirectURL: req.ProviderRedirectURL,
		ExpiresAt:        time.Now().Add(15 * time.Minute),
	})
	if err != nil {
		return "", "", err
	}

	provider, _ := uc.getProvider(req.Provider)
	return provider.BuildAuthURL(stateToken, req.ProviderRedirectURL), req.FinalRedirectURL, nil
}
