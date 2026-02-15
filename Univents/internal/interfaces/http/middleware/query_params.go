package middleware

import (
	"net/http"
	"univents/internal/shared/errx"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/fail/v3"
)

func RequireQueryParams(params ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx, span := MwTracer.Start(ctx, "RequireQueryParams")
			defer span.End()
			q := r.URL.Query()

			for _, p := range params {
				if !q.Has(p) {
					rs, err := fail.ToAs[*resp.Response](fail.New(errx.RequestMissingQueryParam).WithArgs(p).RecordCtx(ctx), "http")
					if err != nil {
						resp.InternalServerError().WithData(err).WithModule("RequireQueryParamsMW").Send(w)
						return
					}
					rs.WithModule("RequireQueryParamsMW").Send(w)
					return
				}
				if q.Get(p) == "" {
					rs, err := fail.ToAs[*resp.Response](fail.New(errx.RequestMissingQueryParamValue).WithArgs(p).RecordCtx(ctx), "http")
					if err != nil {
						resp.InternalServerError().WithData(err).WithModule("RequireQueryParamsMW").Send(w)
						return
					}
					rs.WithModule("RequireQueryParamsMW").Send(w)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequireOnlyQueryParams(allowed ...string) func(http.Handler) http.Handler {
	allowedSet := make(map[string]struct{}, len(allowed))
	for _, p := range allowed {
		allowedSet[p] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx, span := MwTracer.Start(ctx, "RequireOnlyQueryParams")
			defer span.End()
			q := r.URL.Query()

			// check all query params are allowed
			for param := range q {
				if _, ok := allowedSet[param]; !ok {
					rs, err := fail.ToAs[*resp.Response](fail.New(errx.RequestUnknownQueryParam).WithArgs(param).RecordCtx(ctx), "http")
					if err != nil {
						resp.InternalServerError().WithData(err).WithModule("RequireOnlyQueryParamsMW").Send(w)
						return
					}
					rs.WithModule("RequireOnlyQueryParamsMW").Send(w)
					return
				}
			}

			// check all allowed params are present and non-empty
			for _, p := range allowed {
				if !q.Has(p) {
					rs, err := fail.ToAs[*resp.Response](fail.New(errx.RequestMissingQueryParam).WithArgs(p).RecordCtx(ctx), "http")
					if err != nil {
						resp.InternalServerError().WithData(err).WithModule("RequireOnlyQueryParamsMW").Send(w)
						return
					}
					rs.WithModule("RequireOnlyQueryParamsMW").Send(w)
					return
				}
				if q.Get(p) == "" {
					rs, err := fail.ToAs[*resp.Response](fail.New(errx.RequestMissingQueryParamValue).WithArgs(p).RecordCtx(ctx), "http")
					if err != nil {
						resp.InternalServerError().WithData(err).WithModule("RequireOnlyQueryParamsMW").Send(w)
						return
					}
					rs.WithModule("RequireOnlyQueryParamsMW").Send(w)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func DefaultQueryParam(key, value string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			if q.Get(key) == "" {
				q.Set(key, value)
				r.URL.RawQuery = q.Encode()
			}
			next.ServeHTTP(w, r)
		})
	}
}

func AllowOnlyQueryParams(allowed ...string) func(http.Handler) http.Handler {
	allowedSet := make(map[string]struct{}, len(allowed))
	for _, a := range allowed {
		allowedSet[a] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx, span := MwTracer.Start(ctx, "AllowOnlyQueryParams")
			defer span.End()
			q := r.URL.Query()

			for param := range q {
				if _, ok := allowedSet[param]; !ok {
					rs, err := fail.ToAs[*resp.Response](fail.New(errx.RequestUnknownQueryParam).WithArgs(param).RecordCtx(ctx), "http")
					if err != nil {
						resp.InternalServerError().WithData(err).WithModule("AllowOnlyQueryParamsMW").Send(w)
						return
					}
					rs.WithModule("AllowOnlyQueryParamsMW").Send(w)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
