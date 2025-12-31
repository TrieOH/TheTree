package middleware

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/application/authz"
	"GoAuth/internal/domain/auth"
	"GoAuth/internal/ports/outbound"
	"GoAuth/internal/utils"
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type AuthMiddleware struct {
	RevokedRefreshTokensRepo outbound.RevokedRefreshTokenRepository
}

var (
	GoAuthMiddlewareTracer = otel.Tracer("goauth/middleware")
)

func NewAuthMiddleware(RevokedRefreshTokensRepo outbound.RevokedRefreshTokenRepository) *AuthMiddleware {
	return &AuthMiddleware{RevokedRefreshTokensRepo: RevokedRefreshTokensRepo}
}

func (mw *AuthMiddleware) Auth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx, span := GoAuthMiddlewareTracer.Start(ctx, "Middleware.Auth")
			defer span.End()

			var err error
			defer func() {
				if span != nil {
					span.SetAttributes(attribute.Bool("success", err == nil))
				}
			}()

			accessTokenCookie, err := r.Cookie("access_token")
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					mwErr := apierr.ErrUnauthorized.WithMsg("missing access_token cookie").WithID(apierr.AuthMissingAccessCookie).WithCause(err)
					ErrToResp(mwErr).WithModule("AuthMW").Send(w)
					apierr.RecordDomainError(span, mwErr)
					return
				}
				mwErr := apierr.ErrUnauthorized.WithMsg("invalid access_token cookie").WithID(apierr.AuthInvalidAccessCookie).WithCause(err)
				ErrToResp(mwErr).WithModule("AuthMW").Send(w)
				apierr.RecordDomainError(span, mwErr)
				return
			}

			refreshTokenCookie, err := r.Cookie("refresh_token")
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					mwErr := apierr.ErrUnauthorized.WithMsg("missing refresh_token cookie").WithID(apierr.AuthMissingRefreshCookie).WithCause(err)
					ErrToResp(mwErr).WithModule("AuthMW").Send(w)
					apierr.RecordDomainError(span, mwErr)
					return
				}
				mwErr := apierr.ErrUnauthorized.WithMsg("invalid refresh_token cookie").WithID(apierr.AuthInvalidRefreshCookie).WithCause(err)
				ErrToResp(mwErr).WithModule("AuthMW").Send(w)
				apierr.RecordDomainError(span, mwErr)
				return
			}

			var accessToken *auth.AccessClaims
			accessToken, err = utils.ParseAccessToken(accessTokenCookie.Value, utils.GoAuthPublicKey)
			if err != nil {
				ErrToResp(err).WithModule("AuthMW").Send(w)
				apierr.RecordDomainError(span, err)
				return
			}

			span.SetAttributes(
				attribute.String("user.type", accessToken.Sub.UserType),
				attribute.String("user.id", accessToken.Sub.ID.String()),
				attribute.String("user.session_id", accessToken.Sub.SessionID.String()),
			)

			if accessToken.Sub.ProjectID != nil {
				span.SetAttributes(
					attribute.String("user.project_id", accessToken.Sub.ProjectID.String()),
				)
			}

			var refreshToken *auth.RefreshClaims
			refreshToken, err = utils.ParseRefreshToken(refreshTokenCookie.Value, utils.GoAuthPublicKey)
			if err != nil {
				ErrToResp(err).WithModule("AuthMW").Send(w)
				apierr.RecordDomainError(span, err)
				return
			}

			if accessToken.Issuer != "GoAuth" || refreshToken.Issuer != "GoAuth" {
				mwErr := apierr.ErrUnauthorized.WithMsg("invalid issuer").WithID(apierr.TokenInvalidIssuer)
				ErrToResp(mwErr).WithModule("AuthMW").Send(w)
				apierr.RecordDomainError(span, mwErr)
				return
			}

			var refreshTokenJTI uuid.UUID
			refreshTokenJTI, err = uuid.Parse(refreshToken.ID)
			if err != nil {
				mwErr := apierr.ErrUnauthorized.WithMsg("couldn't parse refresh JTI").WithID(apierr.TokenInvalidID).WithCause(err)
				ErrToResp(mwErr).WithModule("AuthMW").Send(w)
				apierr.RecordDomainError(span, mwErr)
				return
			}

			var isRevoked bool
			isRevoked, err = mw.RevokedRefreshTokensRepo.IsRevoked(ctx, refreshTokenJTI)
			if err != nil {
				mwErr := apierr.FromSQLC(err)
				ErrToResp(mwErr).WithModule("AuthMW").Send(w)
				apierr.RecordDomainError(span, mwErr)
				return
			}

			if isRevoked {
				mwErr := apierr.ErrUnauthorized.WithMsg("refresh token is revoked").WithID(apierr.TokenRevoked)
				ErrToResp(mwErr).WithModule("AuthMW").Send(w)
				apierr.RecordDomainError(span, mwErr)
				return
			}

			var principal *authz.Principal
			principal, err = authz.NewPrincipal(accessToken, refreshToken)
			if err != nil {
				ErrToResp(err).WithModule("AuthMW").Send(w)
				apierr.RecordDomainError(span, err)
				return
			}

			ctx = authz.WithPrincipal(ctx, principal)
			ctx = context.WithValue(ctx, userIDKey, principal.UserID.String())

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
