package middleware

import (
	"GoAuth/internal/apierr"
	"net/http"
)

func RequireQueryParams(params ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()

			for _, p := range params {
				if q.Get(p) == "" {
					err := apierr.ErrInvalidInput.WithMsg("missing query parameter: " + p).WithID(apierr.RequestMissingQueryParams)
					ErrToResp(err).WithModule("RequireQueryParamsMW").Send(w)
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
