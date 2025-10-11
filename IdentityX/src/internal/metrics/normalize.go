package metrics

import (
	"net/http"
	"strings"
)

func normalizePath(r *http.Request) string {
	path := r.URL.Path

	return path
}
