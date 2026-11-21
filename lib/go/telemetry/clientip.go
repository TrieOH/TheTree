package telemetry

import (
	"context"
	"net/netip"

	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// MaskIP keeps only the network of a client address: the /24 of an IPv4
// address and the /48 of an IPv6 one. A client IP is personal data, so
// nothing exported (logs, spans) carries a whole one; the network is still
// enough to spot an abusive range or a country. Anything that doesn't parse
// as an address comes back empty.
func MaskIP(s string) string {
	ip, err := netip.ParseAddr(s)
	if err != nil {
		ap, err := netip.ParseAddrPort(s)
		if err != nil {
			return ""
		}
		ip = ap.Addr()
	}
	ip = ip.Unmap()
	bits := 48
	if ip.Is4() {
		bits = 24
	}
	return netip.PrefixFrom(ip, bits).Masked().String()
}

// clientAddrKeys are the span attributes that hold a client address.
// otelhttp sets both at span start.
var clientAddrKeys = []attribute.Key{"client.address", "network.peer.address"}

// ipMasker masks the client address attributes on every span as it starts,
// before any exporter sees it.
type ipMasker struct{}

func (ipMasker) OnStart(_ context.Context, s sdktrace.ReadWriteSpan) {
	for _, kv := range s.Attributes() {
		for _, k := range clientAddrKeys {
			if kv.Key == k {
				s.SetAttributes(attribute.String(string(k), MaskIP(kv.Value.AsString())))
			}
		}
	}
}

func (ipMasker) OnEnd(sdktrace.ReadOnlySpan)      {}
func (ipMasker) Shutdown(context.Context) error   { return nil }
func (ipMasker) ForceFlush(context.Context) error { return nil }
