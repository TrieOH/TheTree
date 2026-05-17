import { Timer } from "../Timer"
import type { ReservedItemI, UnavailableItemI } from "@/features/products/model"

interface CheckoutPartialReservationProps {
  reserved: ReservedItemI[]
  unavailable: UnavailableItemI[]
  confirmDeadline: string
  totalCents: number
  onConfirm: () => void
  onCancel: () => void
}

function formatCurrency(cents: number) {
  return new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency: "BRL"
  }).format(cents / 100)
}

export default function CheckoutPartialReservation({
  reserved,
  unavailable,
  totalCents,
  onConfirm,
  onCancel,
  confirmDeadline
}: CheckoutPartialReservationProps) {
  return (
    <main className="w-full min-w-75 max-w-lg mx-auto px-3 py-8 space-y-6">
      <Timer
        expiresAt={confirmDeadline}
        label="Tempo para confirmar"
      />
      <div>
        <h1 className="text-lg font-bold text-foreground">Reserva parcial</h1>
        <p className="text-sm text-muted-foreground mt-1">
          Alguns itens não estão mais disponíveis. Confira o que foi reservado.
        </p>
      </div>

      {/* Unavailable items */}
      <div className="rounded-md border border-destructive/30 bg-destructive/5 p-4 space-y-2">
        <p className="text-xs font-semibold text-destructive uppercase tracking-wide">
          Itens indisponíveis
        </p>
        {unavailable.map((item) => (
          <div key={item.product_id} className="flex justify-between text-sm">
            <span className="text-foreground">{item.name}</span>
            <span className="text-muted-foreground">{item.reason}</span>
          </div>
        ))}
      </div>

      {/* Reserved items */}
      <div className="rounded-md border border-border p-4 space-y-2">
        <p className="text-xs font-semibold text-muted-foreground uppercase tracking-wide">
          Itens reservados
        </p>
        {reserved.map((item) => (
          <div key={item.product_id} className="flex justify-between text-sm">
            <span className="text-foreground">
              {item.name} × {item.quantity}
            </span>
            <span className="text-foreground font-medium">
              {formatCurrency(item.price_cents * item.quantity)}
            </span>
          </div>
        ))}
        <div className="pt-2 border-t border-border flex justify-between text-sm font-semibold">
          <span>Total</span>
          <span>{formatCurrency(totalCents)}</span>
        </div>
      </div>

      {/* Actions */}
      <div className="flex flex-col gap-2">
        <button
          onClick={onConfirm}
          className="w-full rounded-md bg-primary text-primary-foreground px-4 py-2.5 text-sm font-medium hover:bg-primary/90 transition-colors"
        >
          Continuar com itens disponíveis
        </button>
        <button
          onClick={onCancel}
          className="w-full rounded-md border border-border px-4 py-2.5 text-sm text-muted-foreground hover:bg-muted/50 transition-colors"
        >
          Cancelar
        </button>
      </div>
    </main>
  )
}