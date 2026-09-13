import type { JSX } from "@solidjs/web";
import { createMemo, Show } from "solid-js";
import { Link } from "@tanstack/solid-router";
import BadgeIcon from "~icons/lucide/badge";
import FileCheckIcon from "~icons/lucide/file-check-2";
import ShoppingBagIcon from "~icons/lucide/shopping-bag";

const Badge = BadgeIcon as unknown as () => JSX.Element;
const FileCheck = FileCheckIcon as unknown as () => JSX.Element;
const ShoppingBag = ShoppingBagIcon as unknown as () => JSX.Element;

type Tab = "badges" | "certificates" | "purchases";

export function ProfileCollectionItem(props: { tab: Tab; item: unknown }) {
  const item = createMemo(() =>
    props.item && typeof props.item === "object"
      ? (props.item as Record<string, unknown>)
      : {},
  );
  const title = createMemo(() =>
    String(
      item().template_name ??
        item().name ??
        item().title ??
        item().edition_name ??
        (props.tab === "purchases"
          ? "Pedido"
          : props.tab === "certificates"
            ? "Certificado"
            : "Crachá"),
    ),
  );
  const amount = createMemo(() =>
    typeof item().total_cents === "number"
      ? new Intl.NumberFormat("pt-BR", {
          style: "currency",
          currency: String(item().currency ?? "BRL"),
        }).format((item().total_cents as number) / 100)
      : undefined,
  );
  const description = createMemo(
    () =>
      amount() ??
      item().description ??
      item().status ??
      item().event_name ??
      (typeof item().issued_at === "string"
        ? `Emitido em ${new Date(item().issued_at as string).toLocaleDateString("pt-BR")}`
        : undefined),
  );
  const image = createMemo(() =>
    typeof item().image === "string" ? (item().image as string) : undefined,
  );
  const Icon = () =>
    props.tab === "badges"
      ? Badge
      : props.tab === "certificates"
        ? FileCheck
        : ShoppingBag;

  return (
    <article class="group overflow-hidden rounded-md border border-border bg-card shadow-sm transition hover:-translate-y-0.5 hover:border-primary/50 hover:shadow-md">
      <Show when={image()}>
        {(img) => <img src={img()} alt="" class="h-32 w-full object-cover" />}
      </Show>
      <div class="flex gap-3 p-4">
        <span class="flex size-9 shrink-0 items-center justify-center rounded-md bg-primary/10 text-primary">
          <DynamicIcon component={Icon()} />
        </span>
        <div class="min-w-0">
          <h3 class="truncate text-sm font-semibold text-foreground">
            {title()}
          </h3>
          <Show when={description()}>
            {(desc) => (
              <p class="mt-1 line-clamp-2 text-xs text-muted-foreground">
                {String(desc())}
              </p>
            )}
          </Show>
        </div>
      </div>
      <Show when={props.tab === "purchases" && typeof item().purchase_id === "string"}>
        <Link
          to="/checkouts/$purchaseId"
          params={{ purchaseId: String(item().purchase_id) }}
          class="mx-4 mb-4 flex h-9 items-center justify-center rounded-md border border-border text-sm font-medium hover:bg-muted"
        >
          Ver pedido
        </Link>
      </Show>
    </article>
  );
}

function DynamicIcon(props: { component: () => JSX.Element }) {
  return <>{props.component()}</>;
}
