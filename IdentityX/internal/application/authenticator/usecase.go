package authenticator

import (
	"GoAuth/internal/application/validation"
	"GoAuth/internal/domain/auth"
	"GoAuth/internal/domain/authz"
	"GoAuth/internal/errx"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbounds"
	"context"

	"github.com/MintzyG/fail/v3"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type UseCase struct {
	deps   Deps
	tracer trace.Tracer
}

var _ inbounds.RequestAuthenticator = (*UseCase)(nil)

type Deps struct {
	Session       outbounds.SessionRepository
	Project       outbounds.ProjectRepository
	TokenVerifier inbounds.TokenVerifier
	ApiKey        inbounds.ApiKeyService
}

func New(
	deps Deps,
	tracer trace.Tracer,
) inbounds.RequestAuthenticator {
	return &UseCase{
		deps:   deps,
		tracer: tracer,
	}
}

// AuthenticateRequest
// This function should only be called by AuthMW and therefore does not log errors on the trace
// Leaving this responsibility up to the AuthMW
func (uc *UseCase) AuthenticateRequest(ctx context.Context, in inbounds.AuthenticateRequestInput) (*authz.Principal, error) {
	ctx, span := uc.tracer.Start(ctx, "Authenticator.AuthenticateRequest")
	defer span.End()

	if in.ApiKey != "" {
		span.SetAttributes(attribute.String("auth.method", string(authz.AuthMethodApiKey)))
		return uc.deps.ApiKey.Authenticate(ctx, in.ApiKey)
	}

	span.SetAttributes(attribute.String("auth.method", string(authz.AuthMethodSession)))

	tokenVerifier := uc.deps.TokenVerifier
	sessions := uc.deps.Session

	if in.AccessToken == "" {
		return nil, fail.New(errx.RequestEmptyCookie).WithArgs("access_token").RecordCtx(ctx)
	}
	if in.RefreshToken == "" {
		return nil, fail.New(errx.RequestEmptyCookie).WithArgs("refresh_token").RecordCtx(ctx)
	}

	accessToken, err := tokenVerifier.VerifyAccessToken(ctx, in.AccessToken)
	if err != nil {
		return nil, err
	}

	if accessToken.Sub.ProjectID != nil {
		span.SetAttributes(attribute.String("user.project_id", accessToken.Sub.ProjectID.String()))
	}

	refreshToken, err := tokenVerifier.VerifyRefreshToken(ctx, in.RefreshToken)
	if err != nil {
		return nil, err
	}

	if err = validateIssuers(ctx, in, accessToken, refreshToken); err != nil {
		return nil, err
	}

	refreshTokenJTI, err := validation.RequireRefreshJTI(&refreshToken.ID)
	if err != nil {
		return nil, err
	}

	accessTokenJTI, err := validation.RequireAccessJTI(&accessToken.ID)
	if err != nil {
		return nil, err
	}

	if accessTokenJTI != refreshToken.Sub.AccessJTI {
		return nil, fail.New(errx.TokenMismatchDuringAuth).RecordCtx(ctx)
	}

	sess, err := sessions.GetByFamilyID(ctx, refreshToken.Sub.FamilyID)
	if err != nil {
		if fail.Is(err, errx.SQLNotFound) {
			return nil, fail.New(errx.SessionUnauthorized).RecordCtx(ctx)
		}
		return nil, err
	}

	if sess.SessionID != accessToken.Sub.SessionID {
		return nil, fail.New(errx.TokenSessionMismatch).RecordCtx(ctx)
	}

	// FIXME add occurrence to the audit when its implemented
	if sess.TokenID != refreshTokenJTI {
		_ = sessions.MarkRevokedByFamilyID(ctx, refreshToken.Sub.FamilyID)
		return nil, fail.New(errx.TokenReuseIdentified).WithArgs("refresh").RecordCtx(ctx)
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
	principal, err = authz.NewPrincipal(ctx, accessToken, refreshToken)
	if err != nil {
		return nil, err
	}
	return principal, nil
}

func validateIssuers(
	ctx context.Context,
	in inbounds.AuthenticateRequestInput,
	access *auth.AccessClaims,
	refresh *auth.RefreshClaims,
) error {
	if access.Sub.ProjectID != nil {
		if access.Issuer != access.Sub.ProjectID.String() {
			return fail.New(errx.TokenInvalidIssuer).WithArgs("access").RecordCtx(ctx)
		}
	} else if access.Issuer != in.Issuer {
		return fail.New(errx.TokenInvalidIssuer).WithArgs("access").RecordCtx(ctx)
	}

	if refresh.Issuer != in.Issuer {
		return fail.New(errx.TokenInvalidIssuer).WithArgs("refresh").RecordCtx(ctx)
	}

	return nil
}
