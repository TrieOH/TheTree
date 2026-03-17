import { initMercadoPago } from '@mercadopago/sdk-react'
import { env } from '@/env'

if (typeof window !== 'undefined') {
  initMercadoPago(env.VITE_MERCADO_PAGO_PUBLIC_KEY, {
    locale: 'pt-BR'
  })
}

