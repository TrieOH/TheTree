import { createFileRoute, Link } from "@tanstack/solid-router";
import { useQueryClient } from "@trieoh/front-core/solid";
import { useAuth } from "@trieoh/identityx-sdk-ts-solid";
import type {
  Edition,
  Product,
  ProductVariant,
  StoreStockItem,
  TicketType,
  MyTicket,
} from "@trieoh/univents-api/schemas";
import {
  For,
  Loading,
  Show,
  createMemo,
  createSignal,
  untrack,
} from "solid-js";
import { allPublicEditionsQueryOptions } from "@/features/editions/api";
import { publicEventBySlugQueryOptions } from "@/features/events/api";
import type { EventI } from "@/features/events/model";
import {
  productsQueryOptions,
  productVariantsQueryOptions,
  storeStockQueryOptions,
} from "@/features/products/api";
import { useInventoryStream } from "@/features/products/hooks/use-inventory-stream";
import { EventCart } from "@/features/products/ui/EventCart";
import { ProductCard } from "@/features/products/ui/ProductCard";
import {
  myTicketQueryOptions,
  ticketsQueryOptions,
} from "@/features/tickets/api";
import { TicketCard } from "@/features/tickets/ui/TicketCard";

export const Route = createFileRoute("/events/$slug/store")({
  validateSearch: (search: Record<string, unknown>) => ({
    tab:
      search.tab === "products" ? ("products" as const) : ("tickets" as const),
  }),
  head: ({ params }) => ({
    meta: [{ title: `Loja ${params.slug} - Univents` }],
  }),
  component: StorePage,
});

function StorePage() {
  const params = Route.useParams();
  const search = Route.useSearch();
  const queryClient = useQueryClient();
  const { isAuthenticated } = useAuth();
  const data = createMemo(async () => {
    const event = await queryClient.fetchQuery(
      publicEventBySlugQueryOptions(params().slug),
    );
    if (!event) return null;
    const editions = await queryClient.fetchQuery(
      allPublicEditionsQueryOptions(event.id),
    );
    const edition = currentEdition(editions);
    if (!edition) return { event, edition: null };
    const [tickets, products, stock, heldTicket] = await Promise.all([
      queryClient.fetchQuery(ticketsQueryOptions(edition.id)),
      queryClient.fetchQuery(productsQueryOptions(edition.id)),
      queryClient.fetchQuery(storeStockQueryOptions(edition.id)),
      isAuthenticated()
        ? queryClient.fetchQuery(myTicketQueryOptions(edition.id))
        : null,
    ]);
    const productItems = await Promise.all(
      products.map(async (product) => ({
        product,
        variants: await queryClient.fetchQuery(
          productVariantsQueryOptions(product.id),
        ),
      })),
    );
    return {
      event,
      edition,
      tickets,
      products: productItems,
      stock,
      heldTicket,
    };
  });

  return (
    <Loading fallback={<main class="min-h-screen animate-pulse bg-muted" />}>
      <Show
        when={data()}
        fallback={<main class="p-12 text-center">Evento não encontrado.</main>}
      >
        {(loaded) => (
          <Show
            when={loaded().edition}
            fallback={
              <main class="mx-auto max-w-3xl px-4 py-16 text-center">
                A loja não está disponível no momento.
              </main>
            }
          >
            {(edition) => (
              <StoreContent
                event={loaded().event}
                edition={edition()}
                tickets={loaded().tickets ?? []}
                products={loaded().products ?? []}
                initialStock={loaded().stock ?? []}
                heldTicket={loaded().heldTicket ?? null}
                initialTab={search().tab}
              />
            )}
          </Show>
        )}
      </Show>
    </Loading>
  );
}

function StoreContent(props: {
  event: EventI;
  edition: Edition;
  tickets: TicketType[];
  products: { product: Product; variants: ProductVariant[] }[];
  initialStock: StoreStockItem[];
  initialTab: "tickets" | "products";
  heldTicket: MyTicket | null;
}) {
  const { editionId, initialStock, initialTab, eventSlug } = untrack(() => ({
    editionId: props.edition.id,
    initialStock: props.initialStock,
    initialTab: props.initialTab,
    eventSlug: props.event.slug,
  }));
  const [tab, setTab] = createSignal(initialTab);
  const stock = useInventoryStream(editionId, initialStock);
  const stockById = createMemo(
    () => new Map(stock().map((item) => [item.id, item.stock])),
  );
  return (
    <main class="min-h-screen bg-background px-4 pb-24 pt-4 md:pt-6">
      <div class="mx-auto max-w-6xl">
        <header class="mx-auto max-w-3xl">
          <Link
            to="/events/$slug"
            params={{ slug: eventSlug }}
            class="text-sm font-medium text-muted-foreground"
          >
            ← Evento
          </Link>
          <div class="mt-5 text-center">
            <p class="text-xs font-semibold uppercase tracking-[0.2em] text-primary">
              Loja do evento
            </p>
            <h1 class="mt-3 text-3xl font-bold tracking-tight sm:text-4xl">
              Ingressos e produtos
            </h1>
            <p class="mx-auto mt-3 max-w-xl text-muted-foreground">
              Encontre ingressos e produtos para {props.event.full_name}.
            </p>
          </div>
        </header>
        <div class="mt-10 flex gap-2 border-b border-border/60">
          <For each={["tickets", "products"] as const}>
            {(item) => (
              <button
                type="button"
                class={`border-b-2 px-4 py-3 text-sm font-medium ${tab() === item ? "border-primary text-primary" : "border-transparent text-muted-foreground"}`}
                onClick={() => setTab(item)}
              >
                {item === "tickets" ? "Ingressos" : "Produtos"}
              </button>
            )}
          </For>
        </div>
        <Show
          when={tab() === "tickets"}
          fallback={
            <div class="mt-8 flex flex-wrap justify-center gap-6 sm:justify-start">
              <For
                each={props.products.filter(({ variants }) => variants.length)}
              >
                {({ product, variants }) => (
                  <ProductCard
                    product={product}
                    variants={variants}
                    stock={stockById()}
                    editionId={props.edition.id}
                  />
                )}
              </For>
            </div>
          }
        >
          <div class="mt-8 flex flex-wrap justify-center gap-6 sm:justify-start">
            <For each={props.tickets}>
              {(ticket) => (
                <TicketCard
                  ticket={ticket}
                  editionId={props.edition.id}
                  stock={stockById().get(ticket.id)}
                  heldTicket={props.heldTicket}
                />
              )}
            </For>
          </div>
        </Show>
      </div>
      <EventCart
        editionId={props.edition.id}
        checkoutHref={`/events/${props.event.slug}/checkout`}
      />
    </main>
  );
}

function currentEdition(editions: Edition[]) {
  const now = Date.now();
  const sorted = [...editions].sort((a, b) =>
    a.starts_at.localeCompare(b.starts_at),
  );
  return (
    sorted.find(
      (edition) =>
        new Date(edition.starts_at).getTime() <= now &&
        new Date(edition.ends_at).getTime() >= now,
    ) ??
    sorted.find((edition) => new Date(edition.starts_at).getTime() > now) ??
    sorted.at(-1)
  );
}
