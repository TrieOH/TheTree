// OrderSummary.tsx

import { Package } from "lucide-react";
import { useMemo } from "react";
import type { ReservedItemI } from "@/features/products/model";
import type { CartItem } from "@/features/products/model/cart";

interface OrderSummaryProps {
  items: (CartItem | ReservedItemI)[];
  totalCents?: number;
  title?: string;
  itemCount?: number;
}

interface NormalizedItem {
  id: string;
  name: string;
  quantity: number;
  price_cents: number;
}

function normalizeItem(item: CartItem | ReservedItemI): NormalizedItem {
  return {
    id: "id" in item ? item.id : item.product_id,
    name: item.name,
    quantity: item.quantity,
    price_cents: item.price_cents,
  };
}

function formatCurrency(cents: number) {
  if (cents === 0) return "Gratuito";
  return new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency: "BRL",
  }).format(cents / 100);
}

export function OrderSummary({
  items,
  totalCents: propTotal,
  title = "Resumo",
  itemCount,
}: OrderSummaryProps) {
  const normalizedItems = useMemo(() => items.map(normalizeItem), [items]);

  const total = useMemo(() => {
    if (propTotal !== undefined) return propTotal;
    return normalizedItems.reduce(
      (sum, item) => sum + item.price_cents * item.quantity,
      0,
    );
  }, [propTotal, normalizedItems]);

  const totalItems =
    itemCount ?? normalizedItems.reduce((sum, i) => sum + i.quantity, 0);

  return (
    <div className="w-full min-w-0">
      <div className="flex items-center gap-3 border-b border-border pb-4">
        <span className="flex size-9 shrink-0 items-center justify-center bg-primary/10 text-primary">
          <Package className="size-4" />
        </span>
        <h2 className="font-semibold text-foreground">{title}</h2>
        <span className="ml-auto border border-border bg-muted/40 px-2 py-1 text-xs font-medium text-muted-foreground">
          {totalItems} {totalItems === 1 ? "item" : "itens"}
        </span>
      </div>

      <div className="divide-y divide-border/60">
        {normalizedItems.map((item) => {
          const subtotal = item.price_cents * item.quantity;

          return (
            <div key={item.id} className="flex items-center gap-3 py-4">
              <span className="min-w-6 shrink-0 text-xs font-semibold text-primary">
                {item.quantity}×
              </span>
              <div className="min-w-0 flex-1">
                <p className="text-sm font-semibold leading-snug text-foreground">
                  {item.name}
                </p>
                <p className="text-xs text-muted-foreground">
                  {formatCurrency(item.price_cents)} por unidade
                </p>
              </div>
              <span className="shrink-0 text-sm font-bold tabular-nums text-foreground">
                {formatCurrency(subtotal)}
              </span>
            </div>
          );
        })}
      </div>

      <div className="flex items-end justify-between gap-4 border-t border-border px-1 py-3">
        <div>
          <p className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
            Total do pedido
          </p>
          <p className="mt-0.5 text-xs text-muted-foreground">
            {totalItems}{" "}
            {totalItems === 1 ? "item selecionado" : "itens selecionados"}
          </p>
        </div>
        <span className="text-2xl font-bold text-primary tabular-nums">
          {formatCurrency(total)}
        </span>
      </div>
    </div>
  );
}
