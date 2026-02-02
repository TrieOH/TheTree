package authenticator

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/application/validation"
	"GoAuth/internal/domain/auth"
	"GoAuth/internal/domain/authz"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbounds"
	"context"

	"github.com/MintzyG/fail"
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

	tokenVerifier := uc.deps.TokenVerifier
	sessions := uc.deps.Session

	if in.AccessToken == "" {
		return nil, fail.New(apierr.RequestEmptyCookie).WithArgs("access_token")
	}
	if in.RefreshToken == "" {
		return nil, fail.New(apierr.RequestEmptyCookie).WithArgs("refresh_token")
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

	if err = validateIssuers(in, accessToken, refreshToken); err != nil {
		return nil, err
	}

	refreshTokenJTI, err := validation.RequireRefreshJTI(&refreshToken.ID)
	if err != nil {
		return nil, apierr.FromService(span, err)
	}

	accessTokenJTI, err := validation.RequireAccessJTI(&accessToken.ID)
	if err != nil {
		return nil, apierr.FromService(span, err)
	}

	if accessTokenJTI != refreshToken.Sub.AccessJTI {
		return nil, fail.New(apierr.TokenMismatchDuringAuth)
	}

	sess, err := sessions.GetByFamilyID(ctx, refreshToken.Sub.FamilyID)
	if err != nil {
		if fail.Is(err, apierr.SQLNotFound) {
			return nil, fail.New(apierr.SessionUnauthorized)
		}
		return nil, err
	}

	if sess.SessionID != accessToken.Sub.SessionID {
		return nil, fail.New(apierr.TokenSessionMismatch)
	}

	// FIXME add occurrence to the audit when its implemented
	if sess.TokenID != refreshTokenJTI {
		err = sessions.MarkRevokedByFamilyID(ctx, refreshToken.Sub.FamilyID)
		if err != nil {
			apierr.RecordDomainError(span, err)
		}
		return nil, fail.New(apierr.TokenReuseIdentified).WithArgs("refresh")
	}

	if sess.RevokedAt != nil {
		// should never happen due to query guarding against this, just being defensive
		// system error for appropriate priority if it happens, since it should never happen
		return nil, fail.New(apierr.SessionRevoked)
	}

	span.SetAttributes(
		attribute.String("user.type", accessToken.Sub.UserType),
		attribute.String("user.id", accessToken.Sub.ID.String()),
		attribute.String("user.session_id", accessToken.Sub.SessionID.String()),
	)

	var principal *authz.Principal
	principal, err = authz.NewPrincipal(accessToken, refreshToken)
	if err != nil {
		return nil, apierr.FromService(span, err)
	}
	return principal, nil
}

func validateIssuers(
	in inbounds.AuthenticateRequestInput,
	access *auth.AccessClaims,
	refresh *auth.RefreshClaims,
) error {
	if access.Sub.ProjectID != nil {
		if access.Issuer != access.Sub.ProjectID.String() {
			return fail.New(apierr.TokenInvalidIssuer).WithArgs("access")
		}
	} else if access.Issuer != in.Issuer {
		return fail.New(apierr.TokenInvalidIssuer).WithArgs("access")
	}

	if refresh.Issuer != in.Issuer {
		return fail.New(apierr.TokenInvalidIssuer).WithArgs("refresh")
	}

	return nil
}
