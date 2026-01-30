package middleware

import (
	"GoAuth/internal/apierr"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

func RequireQueryParams(params ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()

			for _, p := range params {
				if !q.Has(p) {
					err := apierr.ErrInvalidInput.WithMsg("missing query parameter: " + p).WithID(apierr.RequestMissingQueryParam)
					resp.FromError(err).WithModule("RequireQueryParamsMW").Send(w)
					return
				}
				if q.Get(p) == "" {
					err := apierr.ErrInvalidInput.WithMsg("missing query parameter value for: " + p).WithID(apierr.RequestMissingQueryParamValue)
					resp.FromError(err).WithModule("RequireQueryParamsMW").Send(w)
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
					err := apierr.ErrInvalidInput.
						WithMsg("unknown query parameter: " + param).
						WithID(apierr.RequestUnknownQueryParam)
					resp.FromError(err).WithModule("RequireOnlyQueryParamsMW").Send(w)
					return
				}
			}

			// check all allowed params are present and non-empty
			for _, p := range allowed {
				if !q.Has(p) {
					err := apierr.ErrInvalidInput.
						WithMsg("missing query parameter: " + p).
						WithID(apierr.RequestMissingQueryParam)
					resp.FromError(err).WithModule("RequireOnlyQueryParamsMW").Send(w)
					return
				}
				if q.Get(p) == "" {
					err := apierr.ErrInvalidInput.
						WithMsg("missing query parameter value for: " + p).
						WithID(apierr.RequestMissingQueryParamValue)
					resp.FromError(err).WithModule("RequireOnlyQueryParamsMW").Send(w)
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
					err := apierr.ErrInvalidInput.
						WithMsg("unknown query parameter: " + key).
						WithID(apierr.RequestUnknownQueryParam)

					resp.FromError(err).
						WithModule("AllowOnlyQueryParamsMW").
						Send(w)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
