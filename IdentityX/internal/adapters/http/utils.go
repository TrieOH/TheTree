package http

import (
	"GoAuth/internal/apierr"
	"errors"
	"net"
	"net/http"
	"strings"
	"time"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

// ErrToResp converts an error to a response.
// It handles API errors and returns a formatted response.
// For unhandled errors, it returns an internal server error response.
func ErrToResp(err error) *resp.Response {
	if err == nil {
		return nil
	}

	var ae *apierr.Error
	if errors.As(err, &ae) {
		return apierr.MapAPIError(ae)
	}

	// unknown error = 500
	return resp.InternalServerError().
		WithTracePrefix("unhandled-error").
		AddTrace(err.Error())
}

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

func CreateCookie(name, value string, age time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   int(time.Until(age).Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
}

func DeleteCookie(name string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
}
