package middleware

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/application/authz"
	"GoAuth/internal/domain/auth"
	"GoAuth/internal/ports/outbound"
	"GoAuth/internal/utils"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type AuthMiddleware struct {
	RevokedRefreshTokensRepo outbound.RevokedRefreshTokenRepository
	tracer                   trace.Tracer
}

func NewAuthMiddleware(RevokedRefreshTokensRepo outbound.RevokedRefreshTokenRepository, tracer trace.Tracer) *AuthMiddleware {
	return &AuthMiddleware{RevokedRefreshTokensRepo: RevokedRefreshTokensRepo, tracer: tracer}
}

// Auth is a middleware function that checks for valid access and refresh tokens.
// It validates the tokens, checks if the refresh token is revoked, and creates a principal from the tokens.
// The principal is then added to the request context.
func (mw *AuthMiddleware) Auth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx, span := mw.tracer.Start(ctx, "Middleware.Auth")
			defer span.End()

			var err error
			defer func() {
				span.SetAttributes(attribute.Bool("success", err == nil))
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
				ErrToResp(err).WithModule("AuthMW").Send(w) // not recording to domain since IsRevoked already does that
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
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
