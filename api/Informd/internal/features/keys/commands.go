package keys

import (
	"Informd/internal/platform/database"
	"Informd/internal/shared/contracts"
	"Informd/internal/shared/ports"
	"context"
	"crypto/rand"
	"encoding/hex"
	authz2 "lib/authz"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
)

type CommandService struct {
	apiKeys  ports.ApiKeysRepo
	projects ports.NamespaceRepo
	perms    authz2.Checker
	tx       database.TxRunner
	tracer   trace.Tracer
}

func NewCommands(
	apiKeys ports.ApiKeysRepo,
	projects ports.NamespaceRepo,
	perms authz2.Checker,
	tx database.TxRunner,
	tracer trace.Tracer,
) *CommandService {
	return &CommandService{
		apiKeys:  apiKeys,
		projects: projects,
		perms:    perms,
		tx:       tx,
		tracer:   tracer,
	}
}

func (s *CommandService) Create(ctx context.Context, keyName string) (rawKey string, ak *contracts.APIKey, err error) {
	ctx, span := s.tracer.Start(ctx, "ApiKeys.Create")
	defer span.End()

	var sub *authz2.UserSubject
	sub, err = authz2.RequireSubject(ctx)
	if err != nil {
		return "", nil, err
	}

	if err = s.perms.Require(ctx,
		authz2.Subject("user", sub.ID),
		authz2.Permission("create_api_key"),
		authz2.Resource("user", sub.ID.String()),
		map[string]any{"subject_id": sub.ID.String()},
	); err != nil {
		return "", nil, err
	}

	rawBytes := make([]byte, 32)
	if _, err := rand.Read(rawBytes); err != nil {
		return "", nil, err
	}
	rawKey = "tf_" + hex.EncodeToString(rawBytes)
	prefix := rawKey[:11] // "tf_" + first 8 hex chars

	hash, err := bcrypt.GenerateFromPassword([]byte(rawKey), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, err
	}

	apiKey, err := contracts.NewAPIKey(sub.ID, keyName, string(hash), prefix)
	if err != nil {
		return "", nil, err
	}

	var created *contracts.APIKey
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

func (s *CommandService) RevokeAPIKey(ctx context.Context, keyID uuid.UUID) error {
	ctx, span := s.tracer.Start(ctx, "ApiKeys.Revoke")
	defer span.End()

	sub, err := authz2.RequireSubject(ctx)
	if err != nil {
		return err
	}

	if err = s.perms.Require(ctx,
		authz2.Subject("user", sub.ID),
		authz2.Permission("revoke"),
		authz2.Resource("api_key", keyID.String()),
		nil,
	); err != nil {
		return err
	}

	if _, err := s.apiKeys.Revoke(ctx, keyID, sub.ID); err != nil {
		return err
	}

	if err = s.perms.DeleteRelation(ctx,
		"api_key:"+keyID.String()+"#parent_user@user:"+sub.ID.String(),
	); err != nil {
		return err
	}

	return nil
}
