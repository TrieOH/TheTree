package logs

import (
	"net/http"
)

func NormalizePath(r *http.Request) string {
	path := r.URL.Path

	return path
}
