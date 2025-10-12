package metrics

import (
	"net/http"
)

func normalizePath(r *http.Request) string {
	path := r.URL.Path

	return path
}
