package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/authz"
	"TriePayments/internal/shared/errx"
	"context"
	"crypto/rand"
	"encoding/hex"
)

func (uc *CommandService) RegisterWebhookEndpoint(ctx context.Context, workspaceName, url string) (*domain.WebhookEndpoint, error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.RegisterWebhookEndpoint")
	defer span.End()

	ga := uc.gaClient

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	workspace, err := uc.workspaces.GetByName(ctx, workspaceName, sub.ID)
	if err != nil {
		return nil, err
	}

	var allowed bool
	allowed, err = ga.Authz.Check().User(sub.ID).
		Object("webhooks").
		Action("create").
		Scope(workspace.ScopeID).
		Allowed(ctx)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, errx.Forbidden("webhooks").SetMessage("insufficient permissions")
	}

	// generate HMAC secret
	secretBytes := make([]byte, 32)
	if _, err := rand.Read(secretBytes); err != nil {
		return nil, err
	}
	secret := hex.EncodeToString(secretBytes)

	endpoint, err := domain.NewWebhookEndpoint(workspace.ID, url, secret)
	if err != nil {
		return nil, err
	}

	created, err := uc.endpoints.Create(ctx, *endpoint)
	if err != nil {
		return nil, err
	}

	return created, nil
}
