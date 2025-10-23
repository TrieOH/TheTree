package metrics

import (
	"net/http"
	"strings"
)

func normalizePath(r *http.Request) string {
    path := r.URL.Path

    switch {
    case strings.HasPrefix(path, "/sessions/"):
        return "/sessions/{session_id}"
    }

    return path
}
