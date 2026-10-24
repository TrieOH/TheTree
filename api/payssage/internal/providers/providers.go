package providers

import (
	"payssage/ports"

	"github.com/MintzyG/fun"
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
		return "", fun.Errf("unsupported provider %s", s).BadRequest()
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
