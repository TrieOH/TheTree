package utils

import (
	"net"
	"net/http"
	"strings"
)

// GetClientIP gets the client's IP address from the request.
// It checks the X-Forwarded-For and X-Real-IP headers first, and then falls back to the remote address.
func GetClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return strings.TrimSpace(strings.Split(ip, ",")[0])
	}

	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return strings.TrimSpace(realIP)
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
