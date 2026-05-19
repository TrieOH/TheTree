package commands

import (
	"Informd/models"
	"context"
	"crypto/rand"
	"encoding/hex"
	"lib/authz"

	"golang.org/x/crypto/bcrypt"
)

func (s *CommandService) Create(ctx context.Context, keyName string) (rawKey string, ak *models.APIKey, err error) {
	ctx, span := s.tracer.Start(ctx, "ApiKeys.Create")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return "", nil, err
	}

	if err = s.perms.Require(ctx,
		authz.Subject("user", sub.ID),
		authz.Permission("create_api_key"),
		authz.Resource("user", sub.ID.String()),
		map[string]any{"subject_id": sub.ID.String()},
	); err != nil {
		return "", nil, err
	}

	rawBytes := make([]byte, 32)
	if _, err = rand.Read(rawBytes); err != nil {
		return "", nil, err
	}
	rawKey = "tf_" + hex.EncodeToString(rawBytes)
	prefix := rawKey[:11] // "tf_" + first 8 hex chars

	hash, err := bcrypt.GenerateFromPassword([]byte(rawKey), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, err
	}

	apiKey, err := models.NewAPIKey(sub.ID, keyName, string(hash), prefix)
	if err != nil {
		return "", nil, err
	}

	var created *models.APIKey
	created, err = s.apiKeys.Create(ctx, *apiKey)
	if err != nil {
		return "", nil, err
	}

	if err = s.perms.CreateRelation(ctx,
		"api_key:"+created.ID.String()+"#parent_user@user:"+sub.ID.String(),
	); err != nil {
		return "", nil, err
	}

	return rawKey, created, nil
}
