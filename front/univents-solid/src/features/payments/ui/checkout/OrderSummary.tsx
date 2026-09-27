import { For, createMemo, type Component } from "solid-js";
import PackageRaw from "~icons/lucide/package";
import type { CartItem } from "@/features/products/model/cart";
import { formatPrice } from "@/shared/lib/money";

const Package = PackageRaw as Component<{ class?: string }>;

interface OrderSummaryProps {
  items: CartItem[];
  totalCents?: number;
  title?: string;
  itemCount?: number;
}

export function OrderSummary(props: OrderSummaryProps) {
  const total = createMemo(() => {
    if (props.totalCents !== undefined) return props.totalCents;
    return props.items.reduce((sum, item) => sum + item.price_cents * item.quantity, 0);
  });

  const totalItems = createMemo(
    () => props.itemCount ?? props.items.reduce((sum, i) => sum + i.quantity, 0),
  );

  return (
    <div class="w-full min-w-0">
      <div class="flex items-center gap-3 border-b border-border pb-4">
        <span class="flex size-9 shrink-0 items-center justify-center bg-primary/10 text-primary">
          <Package class="size-4" />
        </span>
        <h2 class="font-semibold text-foreground">{props.title ?? 'Resumo'}</h2>
        <span class="ml-auto border border-border bg-muted/40 px-2 py-1 text-xs font-medium text-muted-foreground">
          {totalItems()} {totalItems() === 1 ? 'item' : 'itens'}
        </span>
      </div>

      <div class="divide-y divide-border/60">
        <For each={props.items}>
          {(item) => {
            const subtotal = () => item.price_cents * item.quantity;
            return (
              <div class="flex items-center gap-3 py-4">
                <span class="min-w-6 shrink-0 text-xs font-semibold text-primary">{item.quantity}×</span>
                <div class="min-w-0 flex-1">
                  <p class="text-sm font-semibold leading-snug text-foreground">{item.name}</p>
                  <p class="text-xs text-muted-foreground">{formatPrice(item.price_cents)} por unidade</p>
                </div>
                <span class="shrink-0 text-sm font-bold tabular-nums text-foreground">{formatPrice(subtotal())}</span>
              </div>
            );
          }}
        </For>
      </div>

      <div class="flex items-end justify-between gap-4 border-t border-border px-1 py-3">
        <div>
          <p class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">Total do pedido</p>
          <p class="mt-0.5 text-xs text-muted-foreground">
            {totalItems()} {totalItems() === 1 ? 'item selecionado' : 'itens selecionados'}
          </p>
        </div>
        <span class="text-2xl font-bold text-primary tabular-nums">{formatPrice(total())}</span>
      </div>
    </div>
  );
}
