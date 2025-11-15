package utils

import (
	"net/http"
	"strings"
)

var routePatterns = []string{
	"/sessions/{session_id}",
	"/projects/{project_id}",
	"/projects/{project_id}/keys",
}

func NormalizePath(r *http.Request) string {
	path := r.URL.Path
	in := strings.Split(strings.Trim(path, "/"), "/")

	for _, pattern := range routePatterns {
		pat := strings.Split(strings.Trim(pattern, "/"), "/")

		if len(pat) != len(in) {
			continue
		}

		matched := true
		normalized := make([]string, len(pat))

		for i := range pat {
			if strings.HasPrefix(pat[i], "{") && strings.HasSuffix(pat[i], "}") {
				normalized[i] = pat[i]
				continue
			}

			if pat[i] != in[i] {
				matched = false
				break
			}

			normalized[i] = pat[i]
		}

		if matched {
			return "/" + strings.Join(normalized, "/")
		}
	}

	return path
}
