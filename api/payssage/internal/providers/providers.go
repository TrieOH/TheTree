package providers

import (
	"payssage/errs"
	"payssage/ports"
)

type AvailableProviders string

const (
	MercadoPagoProvider AvailableProviders = "mercadopago"
)

func (a AvailableProviders) String() string {
	return string(a)
}

func FromString(s string) (AvailableProviders, error) {
	switch s {
	case "mercadopago":
		return MercadoPagoProvider, nil
	default:
		return "", &errs.UnsupportedProviderError{Provider: s}
	}
}

// PaymentProviders is the registry of provider adapters by provider name.
// Built once at boot (internal/app) and passed to the services that call a
// provider — there is no package-level instance.
type PaymentProviders struct {
	OAuth    map[AvailableProviders]ports.OAuthProvider
	Payments map[AvailableProviders]ports.PaymentAbstractionLayer
	Webhooks map[AvailableProviders]ports.WebhookProvider
}
