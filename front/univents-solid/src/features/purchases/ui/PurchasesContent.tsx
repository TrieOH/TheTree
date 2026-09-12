import { For, Loading, Show, createMemo } from "solid-js";
import type { JSX } from "@solidjs/web";
import type { Purchase } from "@trieoh/univents-api/schemas";
import { Link } from "@tanstack/solid-router";
import { useQueryClient } from "@trieoh/front-core-solid";
import { myPurchasesQueryOptions } from "../api";
import ShoppingBagIcon from "~icons/lucide/shopping-bag";

const ShoppingBag = ShoppingBagIcon as unknown as () => JSX.Element;

const money = (cents: number, currency: string) =>
  new Intl.NumberFormat("pt-BR", { style: "currency", currency }).format(cents / 100);

const status: Record<string, [string, string]> = {
  approved: ["Sucesso", "text-emerald-600 dark:text-emerald-500"],
  pending: ["Pendente", "text-amber-600 dark:text-amber-500"],
  refunded: ["Reembolsado", "text-violet-600 dark:text-violet-500"],
  cancelled: ["Cancelado", "text-muted-foreground"],
  expired: ["Expirado", "text-muted-foreground"],
  failed: ["Falhou", "text-destructive"],
  rejected: ["Recusado", "text-destructive"],
};

export function PurchasesContent() {
  const queryClient = useQueryClient();
  const data = createMemo(() => queryClient.fetchQuery(myPurchasesQueryOptions()));

  return (
    <Loading
      fallback={
        <div class="grid gap-3">
          <div class="h-17 animate-pulse rounded-md bg-muted" />
          <div class="h-17 animate-pulse rounded-md bg-muted" />
        </div>
      }
    >
      <Show when={data()}>
        {(loaded) => (
          <Show
            when={loaded().purchases.length}
            fallback={
              <div class="rounded-md border border-dashed border-border p-10 text-center">
                <h2 class="font-semibold">Você ainda não possui compras</h2>
                <p class="mt-2 text-sm text-muted-foreground">
                  Ingressos, produtos e atividades aparecerão aqui.
                </p>
              </div>
            }
          >
            <div class="grid gap-3 sm:grid-cols-2">
              <For each={loaded().purchases}>
                {(purchase) => <OrderCard purchase={purchase} />}
              </For>
            </div>
          </Show>
        )}
      </Show>
    </Loading>
  );
}

function OrderCard(props: { purchase: Purchase }) {
  const tone = () =>
    status[props.purchase.status] ?? [props.purchase.status, "text-muted-foreground"];

  const date = () =>
    new Intl.DateTimeFormat("pt-BR", { day: "2-digit", month: "2-digit", year: "numeric" }).format(
      new Date(props.purchase.created_at ?? props.purchase.expires_at)
    );

  const payment = () =>
    props.purchase.payment_method === "pix" ? "Pix" : props.purchase.total_cents === 0 ? "" : "Cartão";

  const itemCount = () =>
    `${props.purchase.items.length} ${props.purchase.items.length === 1 ? "item" : "itens"}`;

  const isFree = () => props.purchase.total_cents === 0;

  return (
    <Link
      to="/checkouts/$purchaseId"
      params={{ purchaseId: props.purchase.purchase_id }}
      class="block"
    >
      <article class="flex gap-3 rounded-md border border-border bg-card p-3 shadow-sm transition hover:-translate-y-0.5 hover:border-primary/50 hover:shadow-md sm:p-4">
        <span class="flex size-9 shrink-0 items-center justify-center rounded-md bg-primary/10 text-primary">
          <ShoppingBag />
        </span>

        <div class="min-w-0 flex-1">
          <div class="flex items-center gap-2">
            <h3 class="min-w-0 flex-1 truncate text-sm font-semibold text-foreground">
              {date()}
            </h3>
            <span class={`shrink-0 text-xs font-medium ${tone()[1]}`}>{tone()[0]}</span>
          </div>

          <div class="mt-1 flex items-center justify-between gap-2">
            <p class="min-w-0 truncate text-xs text-muted-foreground">
              {itemCount()}
              {payment() && ` · ${payment()}`}
            </p>
            <Show
              when={!isFree()}
              fallback={<span class="shrink-0 text-xs font-medium text-muted-foreground">Gratuito</span>}
            >
              <strong class="shrink-0 text-sm text-primary">
                {money(props.purchase.total_cents, props.purchase.currency)}
              </strong>
            </Show>
          </div>
        </div>
      </article>
    </Link>
  );
}