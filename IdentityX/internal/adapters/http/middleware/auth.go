package middleware

import (
	"GoAuth/internal/domain/authz"
	"GoAuth/internal/errx"
	"GoAuth/internal/ports/inbounds"
	"errors"
	"net/http"
	"strings"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/fail/v3"
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

			var rs *resp.Response
			var err error
			defer func() {
				span.SetAttributes(attribute.Bool("success", err == nil))
			}()

			var apiKey string
			apiKey = r.Header.Get("X-API-Key")
			if apiKey == "" {
				authHeader := r.Header.Get("Authorization")
				if strings.HasPrefix(authHeader, "Bearer gk_") {
					apiKey = strings.TrimPrefix(authHeader, "Bearer ")
				}
			}

			var in inbounds.AuthenticateRequestInput
			in.Issuer = mw.issuer

			if apiKey != "" {
				in.ApiKey = apiKey
			} else {
				var accessTokenCookie *http.Cookie
				accessTokenCookie, err = r.Cookie("access_token")
				if err != nil {
					if errors.Is(err, http.ErrNoCookie) {
						rs, err = fail.ToAs[*resp.Response](fail.New(errx.AuthMissingAccessCookie).Trace(err.Error()).RecordCtx(ctx), "http")
						if err != nil {
							resp.InternalServerError().WithData(err).WithModule("AuthMW").Send(w)
							return
						}
						rs.WithModule("AuthMW").Send(w)
						return
					}
					rs, err = fail.ToAs[*resp.Response](fail.New(errx.AuthInvalidAccessCookie).Trace(err.Error()).RecordCtx(ctx), "http")
					if err != nil {
						resp.InternalServerError().WithData(err).WithModule("AuthMW").Send(w)
						return
					}
					rs.WithModule("AuthMW").Send(w)
					return
				}

				var refreshTokenCookie *http.Cookie
				refreshTokenCookie, err = r.Cookie("refresh_token")
				if err != nil {
					if errors.Is(err, http.ErrNoCookie) {
						rs, err = fail.ToAs[*resp.Response](fail.New(errx.AuthMissingRefreshCookie).Trace(err.Error()).RecordCtx(ctx), "http")
						if err != nil {
							resp.InternalServerError().WithData(err).WithModule("AuthMW").Send(w)
							return
						}
						rs.WithModule("AuthMW").Send(w)
						return
					}
					rs, err = fail.ToAs[*resp.Response](fail.New(errx.AuthInvalidRefreshCookie).Trace(err.Error()).RecordCtx(ctx), "http")
					if err != nil {
						resp.InternalServerError().WithData(err).WithModule("AuthMW").Send(w)
						return
					}
					rs.WithModule("AuthMW").Send(w)
					return
				}

				in.AccessToken = accessTokenCookie.Value
				in.RefreshToken = refreshTokenCookie.Value
			}

			var principal *authz.Principal
			principal, err = mw.authenticator.AuthenticateRequest(ctx, in)
			if err != nil {
				rs, err = fail.ToAs[*resp.Response](fail.AsFail(err).Trace(err.Error()), "http")
				if err != nil {
					resp.InternalServerError().WithData(err).WithModule("AuthMW").Send(w)
					return
				}
				rs.WithModule("AuthMW").Send(w)
				return
			}

			ctx = authz.WithPrincipal(ctx, principal)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
