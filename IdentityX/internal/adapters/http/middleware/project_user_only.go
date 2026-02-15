package middleware

import (
	"GoAuth/internal/domain/authz"
	"GoAuth/internal/errx"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/fail/v3"
)

func ProjectUserOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, span := MwTracer.Start(ctx, "ProjectUserOnly")
		defer span.End()
		var rs *resp.Response
		principal, err := authz.RequirePrincipal(ctx)
		if err != nil {
			rs, err = fail.ToAs[*resp.Response](fail.AsFail(err).Trace(err.Error()).RecordCtx(ctx), "http")
			if err != nil {
				resp.InternalServerError().WithData(err).WithModule("ProjectUserOnlyMW").Send(w)
				return
			}
			rs.WithModule("ProjectUserOnlyMW").Send(w)
			return
		}

		if principal.ProjectID == nil {
			rs, err = fail.ToAs[*resp.Response](fail.New(errx.AuthNotProjectUser).RecordCtx(ctx), "http")
			if err != nil {
				resp.InternalServerError().WithData(err).WithModule("ProjectUserOnlyMW").Send(w)
				return
			}
			rs.WithModule("ProjectUserOnlyMW").Send(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}
