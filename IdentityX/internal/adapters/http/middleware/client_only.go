package middleware

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/authz"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/fail/v3"
	"go.opentelemetry.io/otel/trace"
)

// ClientOnly is a middleware that ensures that the request is made by a client and not a project user.
func ClientOnly() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx, span := MwTracer.Start(ctx, "ClientOnly")
			trace.ContextWithSpan(ctx, span)
			defer span.End()
			var rs *resp.Response
			principal, err := authz.RequirePrincipal(ctx)
			if err != nil {
				rs, err = fail.ToAs[*resp.Response](fail.AsFail(err).Trace(err.Error()).RecordCtx(ctx), "http")
				if err != nil {
					resp.InternalServerError().WithData(err).WithModule("ClientOnlyMW").Send(w)
					return
				}
				rs.WithModule("ClientOnlyMW").Send(w)
				return
			}

			if principal.ProjectID != nil {
				rs, err = fail.ToAs[*resp.Response](fail.New(apierr.AuthNotClient).RecordCtx(ctx), "http")
				if err != nil {
					resp.InternalServerError().WithData(err).WithModule("ClientOnlyMW").Send(w)
					return
				}
				rs.WithModule("ClientOnlyMW").Send(w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
