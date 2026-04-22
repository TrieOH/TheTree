package api_keys

import (
	"IdentityX/internal/platform/database"
	"IdentityX/internal/shared/authz"
	"IdentityX/internal/shared/contracts"
	"IdentityX/internal/shared/crypto"
	"IdentityX/internal/shared/errx"
	"IdentityX/internal/shared/ports"
	"context"
	"fmt"

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

	principal, err := authz.RequirePrincipal(ctx)
	if err != nil {
		return "", err
	}

	// Only project owner can rotate API key
	isOwner, err := uc.project.IsOwnerOf(ctx, projectID, principal.UserID)
	if err != nil {
		return "", err
	}
	if !isOwner {
		return "", fail.New(errx.ProjectNotFound).RecordCtx(ctx)
	}

	secret, err := crypto.GenerateRandomSecret(32)
	if err != nil {
		return "", fail.New(errx.SYSCryptoError).With(err).RecordCtx(ctx)
	}

	hash, err := crypto.HashBcryptSecret(secret)
	if err != nil {
		return "", fail.New(errx.SYSCryptoError).With(err).RecordCtx(ctx)
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

	principal, err := authz.RequirePrincipal(ctx)
	if err != nil {
		return err
	}

	isOwner, err := uc.project.IsOwnerOf(ctx, projectID, principal.UserID)
	if err != nil {
		return err
	}
	if !isOwner {
		return fail.New(errx.ProjectNotFound).RecordCtx(ctx)
	}

	return uc.apiKeys.Delete(ctx, projectID)
}
