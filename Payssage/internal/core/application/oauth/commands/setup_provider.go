package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/authz"
	"TriePayments/internal/shared/errx"
	"context"
	"time"
)

type SetupProviderRequest struct {
	WorkspaceName    string
	Provider         string
	IsMarketplace    bool
	FeeBps           int
	FinalRedirectURL string
}

func (uc *CommandService) SetupProvider(ctx context.Context, req SetupProviderRequest) (string, error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.SetupProvider")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return "", err
	}

	workspace, err := uc.workspaces.GetByName(ctx, req.WorkspaceName, sub.ID)
	if err != nil {
		return "", err
	}

	if _, err := uc.getProvider(req.Provider); err != nil {
		return "", err
	}

	stateToken, err := generateState()
	if err != nil {
		return "", errx.Internal("oauth_state").SetCause(err)
	}

	_, err = uc.oauthStates.Create(ctx, domain.OAuthState{
		State:            stateToken,
		WorkspaceID:      workspace.ID,
		Provider:         req.Provider,
		Flow:             domain.OAuthFlowSetup,
		IsMarketplace:    req.IsMarketplace,
		FeeBps:           req.FeeBps,
		FinalRedirectURL: req.FinalRedirectURL,
		ExpiresAt:        time.Now().Add(15 * time.Minute),
	})
	if err != nil {
		return "", err
	}

	provider, _ := uc.getProvider(req.Provider)
	return provider.BuildAuthURL(stateToken), nil
}
