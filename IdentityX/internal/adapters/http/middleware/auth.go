package middleware

import (
	"GoAuth/internal/apierr"
	auth2 "GoAuth/internal/application/auth"
	"GoAuth/internal/application/validation"
	"GoAuth/internal/domain/auth"
	authz2 "GoAuth/internal/domain/authz"
	"GoAuth/internal/domain/session"
	authport "GoAuth/internal/ports/auth"
	"GoAuth/internal/ports/outbound"
	"errors"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type AuthMiddleware struct {
	sessions      outbound.SessionRepository
	tokenVerifier authport.TokenVerifier
	tracer        trace.Tracer
	issuer        string
}

func NewAuthMiddleware(
	sessions outbound.SessionRepository,
	tokenVerifier authport.TokenVerifier,
	tracer trace.Tracer,
	issuer string,
) *AuthMiddleware {
	return &AuthMiddleware{
		sessions:      sessions,
		tokenVerifier: tokenVerifier,
		tracer:        tracer,
		issuer:        issuer,
	}
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

			var accessTokenCookie *http.Cookie
			accessTokenCookie, err = r.Cookie("access_token")
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					err = apierr.ErrUnauthorized.WithMsg("missing access_token cookie").WithID(apierr.AuthMissingAccessCookie).WithCause(err)
					resp.FromError(err).WithModule("AuthMW").Send(w)
					apierr.RecordDomainError(span, err)
					return
				}
				err = apierr.ErrUnauthorized.WithMsg("invalid access_token cookie").WithID(apierr.AuthInvalidAccessCookie).WithCause(err)
				resp.FromError(err).WithModule("AuthMW").Send(w)
				apierr.RecordDomainError(span, err)
				return
			}

			var refreshTokenCookie *http.Cookie
			refreshTokenCookie, err = r.Cookie("refresh_token")
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					err = apierr.ErrUnauthorized.WithMsg("missing refresh_token cookie").WithID(apierr.AuthMissingRefreshCookie).WithCause(err)
					resp.FromError(err).WithModule("AuthMW").Send(w)
					apierr.RecordDomainError(span, err)
					return
				}
				err = apierr.ErrUnauthorized.WithMsg("invalid refresh_token cookie").WithID(apierr.AuthInvalidRefreshCookie).WithCause(err)
				resp.FromError(err).WithModule("AuthMW").Send(w)
				apierr.RecordDomainError(span, err)
				return
			}

			var accessToken *auth.AccessClaims
			accessToken, err = mw.tokenVerifier.VerifyAccessToken(ctx, accessTokenCookie.Value)
			if err != nil {
				resp.FromError(err).WithModule("AuthMW").Send(w)
				apierr.RecordDomainError(span, err)
				return
			}

			if accessToken.Sub.ProjectID != nil {
				span.SetAttributes(attribute.String("user.project_id", accessToken.Sub.ProjectID.String()))
			}

			var refreshToken *auth.RefreshClaims
			refreshToken, err = mw.tokenVerifier.VerifyRefreshToken(ctx, refreshTokenCookie.Value)
			if err != nil {
				resp.FromError(err).WithModule("AuthMW").Send(w)
				apierr.RecordDomainError(span, err)
				return
			}

			if accessToken.Issuer != mw.issuer || refreshToken.Issuer != mw.issuer {
				mwErr := apierr.ErrUnauthorized.WithMsg("invalid issuer").WithID(apierr.TokenInvalidIssuer)
				resp.FromError(err).WithModule("AuthMW").Send(w)
				apierr.RecordDomainError(span, mwErr)
				return
			}

			var refreshTokenJTI uuid.UUID
			refreshTokenJTI, err = validation.RequireRefreshJTI(span, &refreshToken.ID)
			if err != nil {
				resp.FromError(err).WithModule("AuthMW").Send(w)
				apierr.RecordDomainError(span, err)
				return
			}

			var accessTokenJTI uuid.UUID
			accessTokenJTI, err = validation.RequireAccessJTI(span, &accessToken.ID)
			if err != nil {
				resp.FromError(err).WithModule("AuthMW").Send(w)
				apierr.RecordDomainError(span, err)
				return
			}

			if accessTokenJTI != refreshToken.Sub.AccessJTI {
				err = apierr.ErrUnauthorized.WithMsg("access token does not belong to this refresh token").WithID(apierr.TokenMismatchDuringAuth)
				apierr.RecordDomainError(span, err)
				resp.FromError(err).WithModule("AuthMW").Send(w)
				return
			}

			var sess *session.Session
			sess, err = mw.sessions.GetByTokenID(ctx, refreshTokenJTI)
			if err != nil {
				if apierr.IsNotFound(err) {
					err = apierr.ErrUnauthorized.WithMsg("session not found or revoked").WithID(apierr.SessionUnauthorized).WithCause(err)
					resp.FromError(err).WithModule("AuthMW").Send(w)
					apierr.RecordDomainError(span, err)
					return
				}
				resp.FromError(err).WithModule("AuthMW").Send(w) // unexpected DB / infra error
				apierr.RecordDomainError(span, err)
				return
			}

			if sess.SessionID != accessToken.Sub.SessionID {
				err = apierr.ErrUnauthorized.WithMsg("token/session mismatch").WithID(apierr.TokenSessionMismatch)
				resp.FromError(err).WithModule("AuthMW").Send(w)
				apierr.RecordDomainError(span, err)
				return
			}

			if sess.RevokedAt != nil {
				// should never happen due to query guarding against this, just being defensive
				// system error for appropriate priority if it happens, since it should never happen
				err = apierr.ErrUnauthorized.WithMsg("session revoked").WithID(apierr.SessionRevoked)
				resp.FromError(err).WithModule("AuthMW").Send(w)
				apierr.RecordSystemError(span, err)
				return
			}

			span.SetAttributes(
				attribute.String("user.type", accessToken.Sub.UserType),
				attribute.String("user.id", accessToken.Sub.ID.String()),
				attribute.String("user.session_id", accessToken.Sub.SessionID.String()),
			)

			var principal *authz2.Principal
			principal, err = authz2.NewPrincipal(accessToken, refreshToken)
			if err != nil {
				resp.FromError(err).WithModule("AuthMW").Send(w)
				apierr.RecordDomainError(span, err)
				return
			}

			ctx = auth2.WithPrincipal(ctx, principal)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
