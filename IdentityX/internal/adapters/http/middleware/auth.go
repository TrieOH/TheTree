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
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type AuthMiddleware struct {
	sessions outbound.SessionRepository
	tracer   trace.Tracer
}

func NewAuthMiddleware(sessions outbound.SessionRepository, tracer trace.Tracer) *AuthMiddleware {
	return &AuthMiddleware{sessions: sessions, tracer: tracer}
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

			issuer := viper.GetString("ISSUER")
			if accessToken.Issuer != issuer || refreshToken.Issuer != issuer {
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

			sess, err := mw.sessions.GetByTokenID(ctx, refreshTokenJTI)
			if err != nil {
				if apierr.IsNotFound(err) {
					mwErr := apierr.ErrUnauthorized.WithMsg("session not found or revoked").WithID(apierr.SessionUnauthorized).WithCause(err)
					ErrToResp(mwErr).WithModule("AuthMW").Send(w)
					apierr.RecordDomainError(span, mwErr)
					return
				}
				ErrToResp(err).WithModule("AuthMW").Send(w) // unexpected DB / infra error
				apierr.RecordDomainError(span, err)
				return
			}

			if sess.SessionID != accessToken.Sub.SessionID {
				mwErr := apierr.ErrUnauthorized.WithMsg("token/session mismatch").WithID(apierr.TokenSessionMismatch)
				ErrToResp(mwErr).WithModule("AuthMW").Send(w)
				apierr.RecordDomainError(span, mwErr)
				return
			}

			if sess.RevokedAt != nil {
				// should never happen due to query guarding against this, just being defensive
				// system error for appropriate priority if it happens, since it should never happen
				mwErr := apierr.ErrUnauthorized.WithMsg("session revoked").WithID(apierr.SessionRevoked)
				ErrToResp(mwErr).WithModule("AuthMW").Send(w)
				apierr.RecordSystemError(span, mwErr)
				return
			}

			span.SetAttributes(
				attribute.String("user.type", accessToken.Sub.UserType),
				attribute.String("user.id", accessToken.Sub.ID.String()),
				attribute.String("user.session_id", accessToken.Sub.SessionID.String()),
			)

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
