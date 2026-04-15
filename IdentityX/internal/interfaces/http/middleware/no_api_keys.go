package middleware

import (
	authz2 "IdentityX/internal/shared/authz"
	"IdentityX/internal/shared/errx"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/fail/v3"
	"go.opentelemetry.io/otel/trace"
)

// NoApiKeys is a middleware that ensures that the request is NOT made using an API key.
func NoApiKeys() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx, span := MwTracer.Start(ctx, "NoApiKeys")
			trace.ContextWithSpan(ctx, span)
			defer span.End()

			var rs *resp.Response
			principal, err := authz2.RequirePrincipal(ctx)
			if err != nil {
				rs, err = fail.ToAs[*resp.Response](fail.AsFail(err).Trace(err.Error()).RecordCtx(ctx), "http")
				if err != nil {
					resp.InternalServerError().WithData(err).WithModule("NoApiKeysMW").Send(w)
					return
				}
				rs.WithModule("NoApiKeysMW").Send(w)
				return
			}

			if principal.Method == authz2.AuthMethodApiKey {
				rs, err = fail.ToAs[*resp.Response](fail.New(errx.AuthApiKeyNotAllowed).RecordCtx(ctx), "http")
				if err != nil {
					resp.InternalServerError().WithData(err).WithModule("NoApiKeysMW").Send(w)
					return
				}
				rs.WithModule("NoApiKeysMW").Send(w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
