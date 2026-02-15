package middleware

import (
	"errors"
	"net/http"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/fail/v3"
	"github.com/TrieOH/goauth-sdk-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type AuthMiddleware struct {
	gaClient goauth.Client
	tracer   trace.Tracer
}

func NewAuthMiddleware(
	gaClient *goauth.Client,
	tracer trace.Tracer,
) *AuthMiddleware {
	return &AuthMiddleware{
		gaClient: *gaClient,
		tracer:   tracer,
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

			accessToken := accessTokenCookie.Value
			token, err := mw.gaClient.Tokens.ValidateToken(ctx, accessToken)
			if err != nil {
				rs, err = fail.ToAs[*resp.Response](fail.AsFail(err), "http")
				if err != nil {
					resp.InternalServerError().WithData(err).WithModule("AuthMW").Send(w)
					return
				}
				rs.WithModule("AuthMW").Send(w)
				return
			}

			subject, err := authz.GetSubjectFromToken(token)
			if err != nil {
				rs, err = fail.ToAs[*resp.Response](fail.AsFail(err), "http")
				if err != nil {
					resp.InternalServerError().WithData(err).WithModule("AuthMW").Send(w)
					return
				}
				rs.WithModule("AuthMW").Send(w)
				return
			}

			ctx = authz.WithSubject(ctx, subject)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
