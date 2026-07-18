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

type PaymentProviders struct {
	OAuth    map[AvailableProviders]ports.OAuthProvider
	Payments map[AvailableProviders]ports.PaymentAbstractionLayer
	Webhooks map[AvailableProviders]ports.WebhookProvider
}

var PayssageProviders PaymentProviders
