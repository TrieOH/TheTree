package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/authz"
	"context"
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func (uc *CommandService) Create(ctx context.Context, workspaceName, keyName string) (rawKey string, ak *domain.APIKey, err error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.Create")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return "", nil, err
	}

	workspace, err := uc.workspaces.GetByName(ctx, workspaceName, sub.ID)
	if err != nil {
		return "", nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("create_api_keys"),
		authz.Resource("workspace", workspace.ID.String()),
	); err != nil {
		return "", nil, err
	}

	rawBytes := make([]byte, 32)
	if _, err := rand.Read(rawBytes); err != nil {
		return "", nil, err
	}
	rawKey = "tp_" + hex.EncodeToString(rawBytes)
	prefix := rawKey[:11] // "tp_" + first 8 hex chars

	hash, err := bcrypt.GenerateFromPassword([]byte(rawKey), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, err
	}

	apiKey, err := domain.NewAPIKey(workspace.ID, keyName, string(hash), prefix)
	if err != nil {
		return "", nil, err
	}

	var created *domain.APIKey
	created, err = uc.apiKeys.Create(ctx, *apiKey)
	if err != nil {
		return "", nil, err
	}

	return rawKey, created, nil
}
