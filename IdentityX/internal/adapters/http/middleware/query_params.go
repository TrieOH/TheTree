package middleware

import (
	"GoAuth/internal/adapters/http/utils"
	"GoAuth/internal/apierr"
	"net/http"

	"github.com/MintzyG/fail"
)

func RequireQueryParams(params ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()

			for _, p := range params {
				if !q.Has(p) {
					rs, ok := utils.Sender(fail.New(apierr.RequestMissingQueryParam).WithArgs(p), "RequireQueryParamsMW", w)
					if ok {
						rs.Send(w)
					}
					return
				}
				if q.Get(p) == "" {
					rs, ok := utils.Sender(fail.New(apierr.RequestMissingQueryParamValue).WithArgs(p), "RequireQueryParamsMW", w)
					if ok {
						rs.Send(w)
					}
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
			q := r.URL.Query()

			// check all query params are allowed
			for param := range q {
				if _, ok := allowedSet[param]; !ok {
					rs, ok := utils.Sender(fail.New(apierr.RequestUnknownQueryParam).WithArgs(param), "RequireOnlyQueryParamsMW", w)
					if ok {
						rs.Send(w)
					}
					return
				}
			}

			// check all allowed params are present and non-empty
			for _, p := range allowed {
				if !q.Has(p) {
					rs, ok := utils.Sender(fail.New(apierr.RequestMissingQueryParam).WithArgs(p), "RequireOnlyQueryParamsMW", w)
					if ok {
						rs.Send(w)
					}
					return
				}
				if q.Get(p) == "" {
					rs, ok := utils.Sender(fail.New(apierr.RequestMissingQueryParamValue).WithArgs(p), "RequireOnlyQueryParamsMW", w)
					if ok {
						rs.WithMsg("missing query parameter value for: " + p).Send(w)
					}
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
			q := r.URL.Query()

			for key := range q {
				if _, ok := allowedSet[key]; !ok {
					rs, ok := utils.Sender(fail.New(apierr.RequestUnknownQueryParam).WithArgs(key), "AllowOnlyQueryParamsMW", w)
					if ok {
						rs.Send(w)
					}
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
