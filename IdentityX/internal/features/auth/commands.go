package auth

import (
	"IdentityX/internal/features/tokens"
	"IdentityX/internal/platform/database"
	"IdentityX/internal/shared/authz"
	"IdentityX/internal/shared/contracts"
	"IdentityX/internal/shared/crypto"
	"IdentityX/internal/shared/errx"
	"IdentityX/internal/shared/ports"
	"context"
	"strings"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type CommandService struct {
	sessions ports.SessionRepository
	project  ports.ProjectRepository
	tokens   tokens.CommandService
	apiKeys  ports.ApiKeyRepository
	logger   *zap.Logger
	tracer   trace.Tracer
	txRunner database.TxRunner
}

func NewCommandService(
	sessions ports.SessionRepository,
	project ports.ProjectRepository,
	tokens tokens.CommandService,
	apiKeys ports.ApiKeyRepository,
	logger *zap.Logger,
	tracer trace.Tracer,
	txRunner database.TxRunner,
) *CommandService {
	return &CommandService{
		sessions: sessions,
		project:  project,
		tokens:   tokens,
		apiKeys:  apiKeys,
		logger:   logger,
		tracer:   tracer,
		txRunner: txRunner,
	}
}

type AuthenticateRequestInput struct {
	AccessToken  string
	RefreshToken string
	ApiKey       string
	Issuer       string
}

// AuthenticateRequest
// This function should only be called by AuthMW and therefore does not log errors on the trace
// Leaving this responsibility up to the AuthMW
func (uc *CommandService) AuthenticateRequest(ctx context.Context, in AuthenticateRequestInput) (*authz.Principal, error) {
	ctx, span := uc.tracer.Start(ctx, "AuthService.AuthenticateRequest")
	defer span.End()

	if in.ApiKey != "" {
		span.SetAttributes(attribute.String("auth.method", string(authz.AuthMethodApiKey)))
		return uc.AuthenticateAPIKey(ctx, in.ApiKey)
	}

	span.SetAttributes(attribute.String("auth.method", string(authz.AuthMethodSession)))
	return uc.AuthenticateSession(ctx, in.AccessToken, in.Issuer)
}

func (uc *CommandService) AuthenticateAPIKey(ctx context.Context, apiKey string) (*authz.Principal, error) {
	ctx, span := uc.tracer.Start(ctx, "AuthService.AuthenticateAPIKey")
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

	keyData, err := uc.apiKeys.GetByProjectID(ctx, projectID)
	if err != nil {
		if errx.IsNotFound(err) {
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
		SessionID: nil,
		Method:    authz.AuthMethodApiKey,
	}, nil
}

func (uc *CommandService) AuthenticateSession(ctx context.Context, accessTokenStr, issuer string) (*authz.Principal, error) {
	ctx, span := uc.tracer.Start(ctx, "AuthService.AuthenticateSession")
	defer span.End()

	if accessTokenStr == "" {
		return nil, fail.New(errx.RequestEmptyCookie).WithArgs("access_token").RecordCtx(ctx)
	}

	accessToken, err := uc.tokens.VerifyAccessToken(ctx, accessTokenStr)
	if err != nil {
		return nil, err
	}

	if accessToken.Sub.ProjectID != nil {
		span.SetAttributes(attribute.String("user.project_id", accessToken.Sub.ProjectID.String()))
	}

	if err = validateIssuers(ctx, issuer, accessToken); err != nil {
		return nil, err
	}

	sess, err := uc.sessions.GetByFamilyID(ctx, accessToken.Sub.FamilyID)
	if err != nil {
		if fail.Is(err, errx.SQLNotFound) {
			return nil, fail.New(errx.SessionUnauthorized).RecordCtx(ctx)
		}
		return nil, err
	}

	if sess.SessionID != accessToken.Sub.SessionID {
		return nil, fail.New(errx.TokenSessionMismatch).RecordCtx(ctx)
	}

	if sess.RevokedAt != nil {
		// should never happen due to query guarding against this, just being defensive
		// system error for appropriate priority if it happens, since it should never happen
		return nil, fail.New(errx.SessionRevoked).RecordCtx(ctx)
	}

	span.SetAttributes(
		attribute.String("user.type", accessToken.Sub.UserType),
		attribute.String("user.id", accessToken.Sub.ID.String()),
		attribute.String("user.session_id", accessToken.Sub.SessionID.String()),
	)

	var principal *authz.Principal
	principal, err = authz.NewPrincipal(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	return principal, nil
}

func validateIssuers(
	ctx context.Context,
	issuer string,
	access *contracts.AccessClaims,
) error {
	if access.Sub.ProjectID != nil {
		if access.Issuer != access.Sub.ProjectID.String() {
			return fail.New(errx.TokenInvalidIssuer).WithArgs("access").RecordCtx(ctx)
		}
	} else if access.Issuer != issuer {
		return fail.New(errx.TokenInvalidIssuer).WithArgs("access").RecordCtx(ctx)
	}

	return nil
}
