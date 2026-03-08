package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/authz"
	"TriePayments/internal/shared/errx"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

type BeginOAuthRequest struct {
	Provider         string
	WorkspaceName    string
	FinalRedirectURL string
}

func (uc *CommandService) BeginOAuth(ctx context.Context, req BeginOAuthRequest) (string, error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.BeginOAuth")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return "", err
	}

	workspace, err := uc.workspaces.GetByName(ctx, req.WorkspaceName, sub.ID)
	if err != nil {
		return "", err
	}

	provider, err := uc.getProvider(req.Provider)
	if err != nil {
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
		FinalRedirectURL: req.FinalRedirectURL,
		ExpiresAt:        time.Now().Add(15 * time.Minute),
	})
	if err != nil {
		return "", err
	}

	redirectURL := provider.BuildAuthURL(stateToken)
	return redirectURL, nil
}

func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (uc *CommandService) getProvider(name string) (domain.OAuthProvider, error) {
	p, ok := uc.providers[name]
	if !ok {
		return nil, errx.Invalid("provider").SetMessage(fmt.Sprintf("unsupported provider: %s", name))
	}
	return p, nil
}
