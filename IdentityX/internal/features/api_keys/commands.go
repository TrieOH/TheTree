package api_keys

import (
	"IdentityX/internal/platform/database"
	authz2 "IdentityX/internal/shared/authz"
	"IdentityX/internal/shared/contracts"
	"IdentityX/internal/shared/crypto"
	errx2 "IdentityX/internal/shared/errx"
	"IdentityX/internal/shared/ports"
	"context"
	"fmt"
	"strings"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type CommandService struct {
	apiKeys  ports.ApiKeyRepository
	project  ports.ProjectRepository
	logger   *zap.Logger
	tracer   trace.Tracer
	txRunner database.TxRunner
}

func NewCommandService(
	apiKeys ports.ApiKeyRepository,
	project ports.ProjectRepository,
	logger *zap.Logger,
	tracer trace.Tracer,
	txRunner database.TxRunner,
) *CommandService {
	return &CommandService{
		apiKeys:  apiKeys,
		project:  project,
		logger:   logger,
		tracer:   tracer,
		txRunner: txRunner,
	}
}

func (uc *CommandService) Rotate(ctx context.Context, projectID uuid.UUID) (string, error) {
	ctx, span := uc.tracer.Start(ctx, "ApiKeyService.Rotate",
		trace.WithAttributes(attribute.String("project.id", projectID.String())),
	)
	defer span.End()

	principal, err := authz2.RequirePrincipal(ctx)
	if err != nil {
		return "", err
	}

	// Only project owner can rotate API key
	isOwner, err := uc.project.IsOwnerOf(ctx, projectID, principal.UserID)
	if err != nil {
		return "", err
	}
	if !isOwner {
		return "", fail.New(errx2.ProjectNotFound).RecordCtx(ctx)
	}

	secret, err := crypto.GenerateRandomSecret(32)
	if err != nil {
		return "", fail.New(errx2.SYSCryptoError).With(err).RecordCtx(ctx)
	}

	hash, err := crypto.HashBcryptSecret(secret)
	if err != nil {
		return "", fail.New(errx2.SYSCryptoError).With(err).RecordCtx(ctx)
	}

	err = uc.apiKeys.Upsert(ctx, contracts.ApiKey{
		ProjectID: projectID,
		ClientID:  principal.UserID,
		KeyHash:   hash,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("gk_%s_%s", projectID.String(), secret), nil
}

func (uc *CommandService) Revoke(ctx context.Context, projectID uuid.UUID) error {
	ctx, span := uc.tracer.Start(ctx, "ApiKeyService.Revoke",
		trace.WithAttributes(attribute.String("project.id", projectID.String())),
	)
	defer span.End()

	principal, err := authz2.RequirePrincipal(ctx)
	if err != nil {
		return err
	}

	isOwner, err := uc.project.IsOwnerOf(ctx, projectID, principal.UserID)
	if err != nil {
		return err
	}
	if !isOwner {
		return fail.New(errx2.ProjectNotFound).RecordCtx(ctx)
	}

	return uc.apiKeys.Delete(ctx, projectID)
}

func (uc *CommandService) Authenticate(ctx context.Context, apiKey string) (*authz2.Principal, error) {
	ctx, span := uc.tracer.Start(ctx, "ApiKeyService.Authenticate")
	defer span.End()

	if !strings.HasPrefix(apiKey, "gk_") {
		return nil, fail.New(errx2.AuthInvalidApiKeyShape).RecordCtx(ctx)
	}

	parts := strings.SplitN(apiKey, "_", 3)
	if len(parts) != 3 {
		return nil, fail.New(errx2.AuthInvalidApiKeyShape).RecordCtx(ctx)
	}

	projectIDStr := parts[1]
	secret := parts[2]

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return nil, fail.New(errx2.AuthInvalidApiKeyShape).RecordCtx(ctx)
	}

	keyData, err := uc.apiKeys.GetByProjectID(ctx, projectID)
	if err != nil {
		if errx2.IsNotFound(err) {
			return nil, fail.New(errx2.AuthInvalidApiKey).RecordCtx(ctx)
		}
		return nil, err
	}

	err = crypto.VerifyBcryptSecret(keyData.KeyHash, secret)
	if err != nil {
		return nil, fail.New(errx2.AuthInvalidApiKey).RecordCtx(ctx)
	}

	return &authz2.Principal{
		UserID:    keyData.ClientID,
		ProjectID: &keyData.ProjectID,
		Method:    authz2.AuthMethodApiKey,
	}, nil
}
