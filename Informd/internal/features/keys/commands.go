package keys

import (
	"TrieForms/internal/plataform/database"
	"TrieForms/internal/shared/authz"
	"TrieForms/internal/shared/ports"
	"TrieForms/internal/shared/types"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"

	fun "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/TrieOH/goauth-sdk-go"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
)

type CommandService struct {
	apiKeys  ports.ApiKeysRepo
	projects ports.ProjectsRepo
	gaClient *goauth.Client
	tx       database.TxRunner
	tracer   trace.Tracer
}

func NewApiKeyCommandService(
	apiKeys ports.ApiKeysRepo,
	projects ports.ProjectsRepo,
	gaClient *goauth.Client,
	tx database.TxRunner,
	tracer trace.Tracer,
) *CommandService {
	return &CommandService{
		apiKeys:  apiKeys,
		projects: projects,
		gaClient: gaClient,
		tx:       tx,
		tracer:   tracer,
	}
}

func (s *CommandService) Create(ctx context.Context, keyName string, projectID uuid.UUID) (rawKey string, ak *types.APIKey, err error) {
	ctx, span := s.tracer.Start(ctx, "ApiKeys.Create")
	defer span.End()

	ga := s.gaClient

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return "", nil, err
	}

	var project *types.Project
	project, err = s.projects.GetByID(ctx, projectID)
	if err != nil {
		return "", nil, err
	}

	var allowed bool
	allowed, err = ga.Authz.Check().User(sub.ID).
		Object("api_keys").
		Action("create").
		Scope(project.ScopeID).
		Allowed(ctx)
	if err != nil {
		return "", nil, err
	}
	if !allowed {
		return "", nil, fun.NewError("insufficient permissions").Forbidden()
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

	apiKey, err := types.NewAPIKey(project.ID, sub.ID, keyName, string(hash), prefix)
	if err != nil {
		return "", nil, err
	}

	meta := json.RawMessage(`{"color": "#de7907", "icon": "KeyRound", "folder": "Api Keys"}`)
	var scope *goauth.Scope
	var idStr = apiKey.ID.String()
	scope, err = ga.Scopes.CreateWithParent(ctx, apiKey.Name, &idStr, &project.ScopeID, meta)
	if err != nil {
		return "", nil, err
	}
	apiKey.AddScope(scope.ID)

	var created *types.APIKey
	created, err = s.apiKeys.Create(ctx, *apiKey)
	if err != nil {
		return "", nil, err
	}

	return rawKey, created, nil
}

func (s *CommandService) RevokeAPIKey(ctx context.Context, projectID, keyID uuid.UUID) error {
	ctx, span := s.tracer.Start(ctx, "ApiKeys.Revoke")
	defer span.End()

	ga := s.gaClient

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	var project *types.Project
	project, err = s.projects.GetByID(ctx, projectID)
	if err != nil {
		return err
	}

	var allowed bool
	allowed, err = ga.Authz.Check().User(sub.ID).
		Object("api_keys").
		Action("revoke").
		Scope(project.ScopeID).
		Allowed(ctx)
	if err != nil {
		return err
	}
	if !allowed {
		return fun.NewError("insufficient permissions").Forbidden()
	}

	if _, err := s.apiKeys.Revoke(ctx, keyID, sub.ID); err != nil {
		return err
	}

	return nil
}
