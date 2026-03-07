package commands

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/shared/authz"
	"TriePayments/internal/shared/errx"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"

	"github.com/TrieOH/goauth-sdk-go"
	"golang.org/x/crypto/bcrypt"
)

func (uc *CommandService) Create(ctx context.Context, workspaceName, keyName string) (rawKey string, ak *domain.APIKey, err error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.Create")
	defer span.End()

	ga := uc.gaClient

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return "", nil, err
	}

	workspace, err := uc.workspaces.GetByName(ctx, workspaceName, sub.ID)
	if err != nil {
		return "", nil, err
	}

	var allowed bool
	allowed, err = ga.Authz.Check().User(sub.ID).
		Object("api_keys").
		Action("create").
		Scope(workspace.ScopeID).
		Allowed(ctx)
	if err != nil {
		return "", nil, err
	}
	if !allowed {
		return "", nil, errx.Forbidden("api key").SetMessage("insufficient permissions")
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

	meta := json.RawMessage(`{"color": "#de7907", "icon": "KeyRound", "folder": "Api Keys"}`)
	var scope *goauth.Scope
	var idStr = apiKey.ID.String()
	scope, err = ga.Scopes.CreateWithParent(ctx, "API Key ("+prefix+")", &idStr, &workspace.ScopeID, meta)
	if err != nil {
		return "", nil, err
	}
	apiKey.AddScope(scope.ID)

	var created *domain.APIKey
	created, err = uc.apiKeys.Create(ctx, *apiKey)
	if err != nil {
		return "", nil, err
	}

	return rawKey, created, nil
}
