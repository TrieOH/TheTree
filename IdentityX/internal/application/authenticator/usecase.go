package authenticator

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/application/validation"
	"GoAuth/internal/domain/authz"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbounds"
	"context"

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

	tokenVerifier := uc.deps.TokenVerifier
	sessions := uc.deps.Session

	if in.AccessToken == "" {
		return nil, apierr.FromService(span, inbounds.ErrEmptyCookie{Cookie: "access_token"})
	}
	if in.RefreshToken == "" {
		return nil, apierr.FromService(span, inbounds.ErrEmptyCookie{Cookie: "refresh_token"})
	}

	accessToken, err := tokenVerifier.VerifyAccessToken(ctx, in.AccessToken)
	if err != nil {
		return nil, apierr.FromService(span, err)
	}

	if accessToken.Sub.ProjectID != nil {
		span.SetAttributes(attribute.String("user.project_id", accessToken.Sub.ProjectID.String()))
	}

	refreshToken, err := tokenVerifier.VerifyRefreshToken(ctx, in.RefreshToken)
	if err != nil {
		return nil, apierr.FromService(span, err)
	}

	// FIXME, later swap this out for checking the issuer of the project if authenticating a project user
	if accessToken.Issuer != in.Issuer {
		return nil, apierr.FromService(span, inbounds.ErrInvalidIssuer{TokenType: "access"})
	}
	if refreshToken.Issuer != in.Issuer {
		return nil, apierr.FromService(span, inbounds.ErrInvalidIssuer{TokenType: "refresh"})
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
		return nil, apierr.FromService(span, inbounds.ErrTokenIDMismatch{})
	}

	sess, err := sessions.GetByTokenID(ctx, refreshTokenJTI)
	if err != nil {
		if apierr.IsNotFound(err) {
			return nil, apierr.FromService(span, inbounds.ErrSessionUnauthorized{})
		}
		return nil, err
	}

	if sess.SessionID != accessToken.Sub.SessionID {
		return nil, apierr.FromService(span, inbounds.ErrTokenSessionMismatch{})
	}

	if sess.RevokedAt != nil {
		// should never happen due to query guarding against this, just being defensive
		// system error for appropriate priority if it happens, since it should never happen
		return nil, apierr.FromService(span, inbounds.ErrAuthSessionRevoked{})
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
