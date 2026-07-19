package providers

import "payssage/internal/features/providers/mercado_pago_provider"

var NewMercadoPago = mercado_pago_provider.NewProvider

type MercadoPagoProvider = mercado_pago_provider.Provider
