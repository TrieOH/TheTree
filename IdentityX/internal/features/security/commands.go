package security

import (
	"IdentityX/internal/platform/database"
	"IdentityX/internal/platform/security"
	"IdentityX/internal/shared/authz"
	"IdentityX/internal/shared/contracts"
	"IdentityX/internal/shared/crypto"
	"IdentityX/internal/shared/errx"
	"IdentityX/internal/shared/feature_deps"
	"IdentityX/internal/shared/ports"
	"context"
	"strings"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type CommandService struct {
	sessions ports.SessionRepository
	project  ports.ProjectRepository
	keys     ports.KeysRepository
	apiKeys  ports.ApiKeyRepository
	logger   *zap.Logger
	tracer   trace.Tracer
	tx       database.TxRunner
}

func NewCommandService(deps feature_deps.SecurityCommandDeps) *CommandService {
	return errx.MustProvide(&CommandService{
		sessions: deps.Sessions,
		project:  deps.Project,
		keys:     deps.Keys,
		apiKeys:  deps.ApiKeys,
		logger:   deps.Logger,
		tracer:   deps.Tracer,
		tx:       deps.Tx,
	})
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
		return nil, fun.ErrUnauthorized("invalid api key shape")
	}

	parts := strings.SplitN(apiKey, "_", 3)
	if len(parts) != 3 {
		return nil, fun.ErrUnauthorized("invalid api key shape")
	}

	projectIDStr := parts[1]
	secret := parts[2]

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return nil, fun.ErrUnauthorized("invalid api key shape")
	}

	keyData, err := uc.apiKeys.GetByProjectID(ctx, projectID)
	if err != nil {
		if fun.Is(err, fun.CodeNotFound) {
			return nil, fun.ErrUnauthorized("invalid api key")
		}
		return nil, err
	}

	err = crypto.VerifyBcryptSecret(keyData.KeyHash, secret)
	if err != nil {
		return nil, fun.ErrUnauthorized("invalid api key")
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
		return nil, fun.ErrUnauthorized("empty access token cookie value")
	}
	accessToken := &contracts.AccessClaims{}
	_, err := security.ParseJWTUnverified[*contracts.AccessClaims](accessTokenStr, accessToken)
	if err != nil {
		return nil, err
	}

	var keyPair *contracts.Pair
	if accessToken.Sub.ProjectID != nil {
		span.SetAttributes(attribute.String("user.project_id", accessToken.Sub.ProjectID.String()))
	}

	keyPair, err = uc.keys.GetActiveSigningKey(ctx, accessToken.Sub.ProjectID)
	if err != nil {
		return nil, err
	}

	accessToken, err = security.VerifyAccessToken(ctx, accessTokenStr, keyPair)
	if err != nil {
		return nil, err
	}

	if err = validateIssuers(issuer, accessToken); err != nil {
		return nil, err
	}

	sess, err := uc.sessions.GetByFamilyID(ctx, accessToken.Sub.FamilyID)
	if err != nil {
		if fun.Is(err, fun.CodeNotFound) {
			return nil, fun.ErrUnauthorized("session not found or revoked")
		}
		return nil, err
	}

	if sess.SessionID != accessToken.Sub.SessionID {
		return nil, fun.ErrUnauthorized("token/session mismatch")
	}
	if sess.RevokedAt != nil {
		// should never happen due to query guarding against this, just being defensive
		// system error for appropriate priority if it happens, since it should never happen
		return nil, fun.ErrUnauthorized("session not found or revoked")
	}

	span.SetAttributes(
		attribute.String("user.type", accessToken.Sub.UserType),
		attribute.String("user.id", accessToken.Sub.ID.String()),
		attribute.String("user.session_id", accessToken.Sub.SessionID.String()),
	)

	var principal *authz.Principal
	principal, err = authz.NewPrincipal(accessToken)
	if err != nil {
		return nil, err
	}
	return principal, nil
}

func validateIssuers(
	issuer string,
	access *contracts.AccessClaims,
) error {
	condition := access.Sub.ProjectID != nil && access.Issuer != access.Sub.ProjectID.String()
	if condition {
		return fun.ErrUnauthorized("access token has invalid issuer")
	}
	if access.Issuer != issuer {
		return fun.ErrUnauthorized("access token has invalid issuer")
	}
	return nil
}
