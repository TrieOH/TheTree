// OrderSummary.tsx
import { useMemo } from "react"
import { Package } from "lucide-react"
import type { CartItem } from "@/features/products/model/cart"
import type { ReservedItemI } from "@/features/products/model"

interface OrderSummaryProps {
  items: (CartItem | ReservedItemI)[]
  totalCents?: number
  title?: string
  itemCount?: number
}

interface NormalizedItem {
  id: string
  name: string
  quantity: number
  price_cents: number
}

function normalizeItem(item: CartItem | ReservedItemI): NormalizedItem {
  return {
    id: "id" in item ? item.id : item.product_id,
    name: item.name,
    quantity: item.quantity,
    price_cents: item.price_cents
  }
}

function formatCurrency(cents: number) {
  return new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency: "BRL"
  }).format(cents / 100)
}

export function OrderSummary({
  items,
  totalCents: propTotal,
  title = "Resumo",
  itemCount
}: OrderSummaryProps) {

  const normalizedItems = useMemo(
    () => items.map(normalizeItem),
    [items]
  )

  const total = useMemo(() => {
    if (propTotal !== undefined) return propTotal
    return normalizedItems.reduce(
      (sum, item) => sum + item.price_cents * item.quantity,
      0
    )
  }, [propTotal, normalizedItems])

  const totalItems = itemCount ?? normalizedItems.reduce((sum, i) => sum + i.quantity, 0)

  return (
    <div className="w-full min-w-75 space-y-3">
      <div className="flex items-center gap-2 pb-2 border-b border-border">
        <Package className="w-4 h-4 text-primary" />
        <h2 className="text-xs font-bold uppercase tracking-wide text-muted-foreground">
          {title}
        </h2>
        <span className="ml-auto text-xs text-muted-foreground">
          {totalItems} itens
        </span>
      </div>

      <div className="space-y-0 divide-y divide-border/50">
        {normalizedItems.map((item) => {
          const subtotal = item.price_cents * item.quantity

          return (
            <div key={item.id} className="flex items-center gap-3 py-2.5">
              <span className="text-xs font-semibold text-primary min-w-6">
                {item.quantity}×
              </span>
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium text-foreground truncate">
                  {item.name}
                </p>
                <p className="text-xs text-muted-foreground">
                  {formatCurrency(item.price_cents)} un
                </p>
              </div>
              <span className="text-sm font-semibold text-foreground tabular-nums">
                {formatCurrency(subtotal)}
              </span>
            </div>
          )
        })}
      </div>

      <div className="flex items-center justify-between pt-3 border-t-2 border-primary/20">
        <span className="text-sm font-bold uppercase tracking-wide text-foreground">
          Total
        </span>
        <span className="text-xl font-bold text-primary tabular-nums">
          {formatCurrency(total)}
        </span>
      </div>
    </div>
  )
}