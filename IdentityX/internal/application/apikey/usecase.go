package apikey

import (
	"GoAuth/internal/crypto"
	"GoAuth/internal/domain/apikey"
	"GoAuth/internal/domain/authz"
	"GoAuth/internal/errx"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbounds"
	"context"
	"fmt"
	"strings"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	usecaseTracer = otel.Tracer("apikey_usecase")
)

type UseCase struct {
	deps Deps
}

type Deps struct {
	ApiKey  outbounds.ApiKeyRepository
	Project outbounds.ProjectRepository
}

var _ inbounds.ApiKeyService = (*UseCase)(nil)

func New(deps Deps) inbounds.ApiKeyService {
	return &UseCase{
		deps: deps,
	}
}

func (uc *UseCase) Rotate(ctx context.Context, projectID uuid.UUID) (string, error) {
	ctx, span := usecaseTracer.Start(ctx, "ApiKeyService.Rotate",
		trace.WithAttributes(attribute.String("project.id", projectID.String())),
	)
	defer span.End()

	principal, err := authz.RequirePrincipal(ctx)
	if err != nil {
		return "", err
	}

	// Only project owner can rotate API key
	isOwner, err := uc.deps.Project.IsOwnerOf(ctx, projectID, principal.UserID)
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

	err = uc.deps.ApiKey.Upsert(ctx, apikey.ApiKey{
		ProjectID: projectID,
		ClientID:  principal.UserID,
		KeyHash:   hash,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("gk_%s_%s", projectID.String(), secret), nil
}

func (uc *UseCase) Revoke(ctx context.Context, projectID uuid.UUID) error {
	ctx, span := usecaseTracer.Start(ctx, "ApiKeyService.Revoke",
		trace.WithAttributes(attribute.String("project.id", projectID.String())),
	)
	defer span.End()

	principal, err := authz.RequirePrincipal(ctx)
	if err != nil {
		return err
	}

	isOwner, err := uc.deps.Project.IsOwnerOf(ctx, projectID, principal.UserID)
	if err != nil {
		return err
	}
	if !isOwner {
		return fail.New(errx.ProjectNotFound).RecordCtx(ctx)
	}

	return uc.deps.ApiKey.Delete(ctx, projectID)
}

func (uc *UseCase) Authenticate(ctx context.Context, apiKey string) (*authz.Principal, error) {
	ctx, span := usecaseTracer.Start(ctx, "ApiKeyService.Authenticate")
	defer span.End()

	if !strings.HasPrefix(apiKey, "gk_") {
		return nil, fail.New(errx.AuthInvalidApiKeyShape).RecordCtx(ctx)
	}

	parts := strings.SplitN(apiKey, "_", 3)
	if len(parts) != 3 {
		return nil, fail.New(errx.AuthInvalidApiKeyShape).RecordCtx(ctx)
	}

	projectIDStr := parts[1]
	secret := parts[2]

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return nil, fail.New(errx.AuthInvalidApiKeyShape).RecordCtx(ctx)
	}

	keyData, err := uc.deps.ApiKey.GetByProjectID(ctx, projectID)
	if err != nil {
		if errx.IsNotFoundNew(err) {
			return nil, fail.New(errx.AuthInvalidApiKey).RecordCtx(ctx)
		}
		return nil, err
	}

	err = crypto.VerifyBcryptSecret(keyData.KeyHash, secret)
	if err != nil {
		return nil, fail.New(errx.AuthInvalidApiKey).RecordCtx(ctx)
	}

	return &authz.Principal{
		UserID:    keyData.ClientID,
		ProjectID: &keyData.ProjectID,
		Method:    authz.AuthMethodApiKey,
	}, nil
}
