package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/authz"
	"TriePayments/internal/shared/errx"
	"context"
	"fmt"
	"time"
)

type ConnectSellerRequest struct {
	WorkspaceName    string
	Provider         string
	FinalRedirectURL string
}

func (uc *CommandService) ConnectSeller(ctx context.Context, req ConnectSellerRequest) (string, error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.ConnectSeller")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return "", err
	}

	workspace, err := uc.workspaces.GetByName(ctx, req.WorkspaceName, sub.ID)
	if err != nil {
		return "", err
	}

	// verify workspace has this provider set up
	creds, err := uc.credentials.ListByWorkspace(ctx, workspace.ID)
	if err != nil {
		return "", err
	}
	hasProvider := false
	for _, c := range creds {
		if c.Provider == req.Provider && c.RevokedAt == nil {
			hasProvider = true
			break
		}
	}
	if !hasProvider {
		return "", errx.Invalid("provider").SetMessage(fmt.Sprintf("workspace has not set up provider: %s", req.Provider))
	}

	stateToken, err := generateState()
	if err != nil {
		return "", errx.Internal("oauth_state").SetCause(err)
	}

	_, err = uc.oauthStates.Create(ctx, domain.OAuthState{
		State:            stateToken,
		WorkspaceID:      workspace.ID,
		Provider:         req.Provider,
		Flow:             domain.OAuthFlowConnect,
		IsMarketplace:    false,
		FeeBps:           0,
		FinalRedirectURL: req.FinalRedirectURL,
		ExpiresAt:        time.Now().Add(15 * time.Minute),
	})
	if err != nil {
		return "", err
	}

	provider, _ := uc.getProvider(req.Provider)
	return provider.BuildAuthURL(stateToken, req.FinalRedirectURL), nil
}
