import type { SubmitPaymentPayloadI } from "@/features/payments/model"
import type { CheckoutPhase } from "@/features/products/hooks/use-checkout-socket"
import type { ReservedItemI } from "@/features/products/model"
import { OrderSummary } from "@/features/payments/ui/checkout/OrderSummary"
import { Timer } from "@/features/payments/ui/checkout/Timer"
import { PaymentProviderSelector } from "@/features/payments/ui/PaymentProviderSelector"

interface CheckoutPaymentFormProps {
  phase: Extract<CheckoutPhase, "reservation_confirmed" | "awaiting_payment" | "payment_processing" | "payment_failed">
  reservedItems: ReservedItemI[]
  totalCents: number
  reservationExpiresAt: string | null
  paymentIntentId: string | null
  sellerPublicKey: string
  onSubmit: (data: SubmitPaymentPayloadI) => void
  onCancel: () => void
  onExpire: () => void
}

export default function CheckoutPaymentForm({
  phase,
  reservedItems,
  totalCents,
  reservationExpiresAt,
  sellerPublicKey,
  onSubmit,
  onCancel,
  onExpire,
}: CheckoutPaymentFormProps) {
  const isProcessing = phase === "payment_processing" || phase === "awaiting_payment"

  const cartItems = reservedItems.map((item) => ({
    id: item.product_id,
    name: item.name,
    price_cents: item.price_cents,
    quantity: item.quantity,
    inventory_remaining: item.quantity,
    has_inventory: true,
  }))

  return (
    <main className="w-full min-w-75 max-w-4xl mx-auto px-3 py-4 space-y-2">
      {/* Timer */}
      {reservationExpiresAt && (
        <Timer
          expiresAt={reservationExpiresAt}
          onExpire={onExpire}
        />
      )}

      {/* Header */}
      <div className="pb-3 mb-4 border-b border-border">
        <h1 className="text-lg font-bold text-foreground">Checkout</h1>
        <p className="text-xs text-muted-foreground">Finalize sua compra</p>
      </div>

      {/* Processing overlay banner */}
      {isProcessing && (
        <div className="flex items-center gap-2 rounded-md bg-muted/60 border border-border px-4 py-3 text-sm text-muted-foreground">
          <div className="w-4 h-4 rounded-full border-2 border-primary border-t-transparent animate-spin shrink-0" />
          <span>Processando pagamento…</span>
        </div>
      )}

      {/* Grid — OrderSummary and PaymentProviderSelector*/}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <OrderSummary items={cartItems} totalCents={totalCents} />

        <div className="lg:pl-6 lg:border-l lg:border-border">
          <PaymentProviderSelector
            amount={totalCents}
            handleSubmit={onSubmit}
            onCancel={onCancel}
            sellerPublicKey={sellerPublicKey}
          />
        </div>
      </div>
    </main>
  )
}