package http

import (
	"GoAuth/internal/adapters/http/validation"
	"GoAuth/internal/apierr"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"strings"
	"time"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func getUUID(r *http.Request, fieldName string) (uuid.UUID, *resp.Response) {
	projectIDStr := chi.URLParam(r, fieldName)
	if err := validation.ValidateRule(projectIDStr, "required,uuid7", fieldName); err != nil {
		return uuid.Nil, resp.FromError(err)
	}

	projectID, err := validation.ParseUUID(projectIDStr, fieldName)
	if err != nil {
		return uuid.Nil, resp.FromError(apierr.FromHandler(err))
	}
	return projectID, nil
}

func SetTrustProxyHeaders() {
	HTTPProxyConfig.TrustProxyHeaders = viper.GetBool("TRUST_PROXY_HEADERS")
}

func SetTrustedProxies() error {
	raw := viper.GetStringSlice("TRUSTED_PROXIES")
	if len(raw) == 0 {
		HTTPProxyConfig.TrustedProxies = nil
		return nil
	}

	proxies := make([]netip.Prefix, 0, len(raw))

	for _, cidr := range raw {
		cidr = strings.TrimSpace(cidr)
		if cidr == "" {
			continue
		}

		prefix, err := netip.ParsePrefix(cidr)
		if err != nil {
			return fmt.Errorf("invalid TRUSTED_PROXIES entry %q: %w", cidr, err)
		}

		proxies = append(proxies, prefix)
	}

	HTTPProxyConfig.TrustedProxies = proxies
	return nil
}

func LoadProxyConfig() error {
	SetTrustProxyHeaders()
	if err := SetTrustedProxies(); err != nil {
		return err
	}
	if HTTPProxyConfig.TrustProxyHeaders && len(HTTPProxyConfig.TrustedProxies) == 0 {
		return errors.New("TRUST_PROXY_HEADERS=true but TRUSTED_PROXIES is empty")
	}
	return nil
}

type ProxyConfig struct {
	// Enable reading X-Forwarded-For / X-Real-IP
	TrustProxyHeaders bool
	TrustedProxies    []netip.Prefix
}

var HTTPProxyConfig = ProxyConfig{}

func ClientIPString(ip netip.Addr) string {
	if !ip.IsValid() {
		return "unknown"
	}
	return ip.String()
}

// GetClientIP gets the client's IP address from the request.
// It checks the X-Forwarded-For and X-Real-IP headers first, and then falls back to the remote address.
func GetClientIP(r *http.Request, cfg ProxyConfig) netip.Addr {
	remoteIP, ok := extractRemoteAddr(r)
	if !ok {
		return netip.Addr{}
	}

	// If we do not trust proxy headers, stop here
	if !cfg.TrustProxyHeaders {
		return remoteIP
	}

	// If the immediate sender is NOT a trusted proxy, do not trust headers
	if !isTrustedProxy(remoteIP, cfg.TrustedProxies) {
		return remoteIP
	}

	// Try X-Forwarded-For first (standard)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if ip, ok := parseXForwardedFor(xff, cfg.TrustedProxies); ok {
			return ip
		}
	}

	// Fallback to X-Real-IP
	if rip := r.Header.Get("X-Real-IP"); rip != "" {
		if ip, err := netip.ParseAddr(strings.TrimSpace(rip)); err == nil {
			return ip
		}
	}

	return remoteIP
}

func extractRemoteAddr(r *http.Request) (netip.Addr, bool) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return netip.Addr{}, false
	}

	ip, err := netip.ParseAddr(host)
	if err != nil {
		return netip.Addr{}, false
	}

	return ip, true
}

func isTrustedProxy(ip netip.Addr, trusted []netip.Prefix) bool {
	for _, p := range trusted {
		if p.Contains(ip) {
			return true
		}
	}
	return false
}

func parseXForwardedFor(xff string, trusted []netip.Prefix) (netip.Addr, bool) {
	parts := strings.Split(xff, ",")

	for i := len(parts) - 1; i >= 0; i-- {
		part := parts[i]
		part = strings.TrimSpace(part)

		ip, err := netip.ParseAddr(part)
		if err != nil {
			continue
		}

		if !isTrustedProxy(ip, trusted) {
			return ip, true
		}
	}

	return netip.Addr{}, false
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
