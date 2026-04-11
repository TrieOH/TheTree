package keys

import (
	"TrieForms/internal/platform/database"
	"TrieForms/internal/shared/authz"
	"TrieForms/internal/shared/contracts"
	"TrieForms/internal/shared/ports"
	"context"
	"crypto/rand"
	"encoding/hex"

	v1 "github.com/authzed/authzed-go/v1"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
)

type CommandService struct {
	apiKeys  ports.ApiKeysRepo
	projects ports.ProjectsRepo
	az       *v1.Client
	tx       database.TxRunner
	tracer   trace.Tracer
}

func NewApiKeyCommandService(
	apiKeys ports.ApiKeysRepo,
	projects ports.ProjectsRepo,
	az *v1.Client,
	tx database.TxRunner,
	tracer trace.Tracer,
) *CommandService {
	return &CommandService{
		apiKeys:  apiKeys,
		projects: projects,
		az:       az,
		tx:       tx,
		tracer:   tracer,
	}
}

func (s *CommandService) Create(ctx context.Context, keyName string, projectID uuid.UUID) (rawKey string, ak *contracts.APIKey, err error) {
	ctx, span := s.tracer.Start(ctx, "ApiKeys.Create")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return "", nil, err
	}

	var project *contracts.Project
	project, err = s.projects.GetByID(ctx, projectID)
	if err != nil {
		return "", nil, err
	}

	if err = authz.Require(ctx, s.az,
		authz.Subject("user", sub.ID),
		authz.Permission("create_key"),
		authz.Resource("project", project.ID.String()),
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

	apiKey, err := contracts.NewAPIKey(project.ID, sub.ID, keyName, string(hash), prefix)
	if err != nil {
		return "", nil, err
	}

	var created *contracts.APIKey
	created, err = s.apiKeys.Create(ctx, *apiKey)
	if err != nil {
		return "", nil, err
	}

	return rawKey, created, nil
}

func (s *CommandService) RevokeAPIKey(ctx context.Context, keyID uuid.UUID) error {
	ctx, span := s.tracer.Start(ctx, "ApiKeys.Revoke")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	if err = authz.Require(ctx, s.az,
		authz.Subject("user", sub.ID),
		authz.Permission("revoke"),
		authz.Resource("api_key", keyID.String()),
	); err != nil {
		return err
	}

	if _, err := s.apiKeys.Revoke(ctx, keyID, sub.ID); err != nil {
		return err
	}

	return nil
}
