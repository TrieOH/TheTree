package providers

import (
	"fmt"
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
		return "", fmt.Errorf("unsupported provider %s", s)
	}
}

type PaymentProviders struct {
	OAuth    map[AvailableProviders]ports.OAuthProvider
	Payments map[AvailableProviders]ports.PaymentAbstractionLayer
}

var PayssageProviders PaymentProviders
