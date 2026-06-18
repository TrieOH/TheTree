package api_keys

import (
	"IdentityX/models"
	"context"
	"crypto/rand"
	"encoding/hex"
	"payssage/ports"

	"lib/authz"
	"lib/database"
	"payssage/models"

	"github.com/authzed/authzed-go/v1"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
)

type CommandService struct {
	apiKeys    ports.ApiKeysRepo
	workspaces ports.WorkspaceRepo
	az         *authzed.Client
	tx         database.TxRunner
	tracer     trace.Tracer
}

func NewCommandService(
	apiKeys ports.ApiKeysRepo,
	workspaces ports.WorkspaceRepo,
	az *authzed.Client,
	tx database.TxRunner,
	tracer trace.Tracer,
) *CommandService {
	return &CommandService{
		apiKeys:    apiKeys,
		workspaces: workspaces,
		az:         az,
		tx:         tx,
		tracer:     tracer,
	}
}

func (uc *CommandService) Create(ctx context.Context, workspaceName, keyName string) (rawKey string, ak *models.APIKey, err error) {
	ctx, span := uc.tracer.Start(ctx, "CommandService.Create")
	defer span.End()

	var sub *models.UserSubject
	sub, err = models.RequireSubject(ctx)
	if err != nil {
		return "", nil, err
	}

	workspace, err := uc.workspaces.GetByName(ctx, workspaceName, sub.ID)
	if err != nil {
		return "", nil, err
	}

	//if err = uc.checker.Require(ctx,
	//	authz.Subject("user", sub.ID),
	//	authz.Permission("create_api_keys"),
	//	authz.Resource("workspace", workspace.ID.String()),
	//	nil,
	//); err != nil {
	//	return "", nil, err
	//}

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

	apiKey, err := models.NewAPIKey(workspace.ID, keyName, string(hash), prefix)
	if err != nil {
		return "", nil, err
	}

	var created *models.APIKey
	created, err = uc.apiKeys.Create(ctx, *apiKey)
	if err != nil {
		return "", nil, err
	}

	return rawKey, created, nil
}

func (uc *CommandService) RevokeAPIKey(ctx context.Context, workspaceName string, keyID uuid.UUID) error {
	ctx, span := uc.tracer.Start(ctx, "CommandService.RevokeAPIKey")
	defer span.End()

	sub, err := models.RequireSubject(ctx)
	if err != nil {
		return err
	}

	workspace, err := uc.workspaces.GetByName(ctx, workspaceName, sub.ID)
	if err != nil {
		return err
	}

	//if err = authz.Require(ctx, uc.az,
	//	authz.Subject("user", sub.ID),
	//	authz.Permission("revoke_api_keys"),
	//	authz.Resource("workspace", workspace.ID.String()),
	//); err != nil {
	//	return err
	//}

	if _, err := uc.apiKeys.Revoke(ctx, keyID, workspace.ID); err != nil {
		return err
	}

	return nil
}
