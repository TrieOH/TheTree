package validation

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"strings"

	"github.com/spf13/viper"
)

func SetTrustProxyHeaders() {
	HTTPProxyConfig.TrustProxyHeaders = viper.GetBool("TRUST_PROXY_HEADERS")
}

func SetTrustedProxies() error {
	raw := viper.GetString("TRUSTED_PROXIES")
	if raw == "" {
		HTTPProxyConfig.TrustedProxies = nil
		return nil
	}

	parts := strings.Split(raw, ",")
	proxies := make([]netip.Prefix, 0, len(parts))

	for _, cidr := range parts {
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
