package middleware

import (
	"GoAuth/internal/apierr"
	appauth "GoAuth/internal/application/auth"
	"GoAuth/internal/domain/authz"
	"GoAuth/internal/ports/inbounds"
	"errors"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type AuthMiddleware struct {
	authenticator inbounds.RequestAuthenticator
	tracer        trace.Tracer
	issuer        string
}

func NewAuthMiddleware(
	authenticator inbounds.RequestAuthenticator,
	tracer trace.Tracer,
	issuer string,
) *AuthMiddleware {
	return &AuthMiddleware{
		authenticator: authenticator,
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

			in := inbounds.AuthenticateRequestInput{
				AccessToken:  accessTokenCookie.Value,
				RefreshToken: refreshTokenCookie.Value,
				Issuer:       mw.issuer,
			}

			var principal *authz.Principal
			principal, err = mw.authenticator.AuthenticateRequest(ctx, in)
			if err != nil {
				resp.FromError(err).WithModule("AuthMW").Send(w)
				return
			}

			ctx = appauth.WithPrincipal(ctx, principal)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
