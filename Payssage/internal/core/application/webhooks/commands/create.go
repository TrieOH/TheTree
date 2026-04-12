package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/authz"
	"context"
	"crypto/rand"
	"encoding/hex"
)

func (uc *CommandService) RegisterWebhookEndpoint(ctx context.Context, workspaceName, url string) (*domain.WebhookEndpoint, error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.RegisterWebhookEndpoint")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	workspace, err := uc.workspaces.GetByName(ctx, workspaceName, sub.ID)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("create_webhooks"),
		authz.Resource("workspace", workspace.ID.String()),
	); err != nil {
		return nil, err
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
