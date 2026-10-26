import { createFileRoute, Link } from "@tanstack/solid-router";
import { useQuery, useQueryClient } from "@trieoh/front-core-solid";
import { useAuth } from "@trieoh/identityx-sdk-ts-solid";
import type {
  Edition,
  Product,
  ProductVariant,
  StoreStockItem,
  TicketType,
  MyTicket,
} from "@trieoh/univents-api/schemas";
import { For, Show, createMemo, createSignal, untrack } from "solid-js";
import type { EventI } from "@/features/events/model";
import { storePageQueryOptions } from "@/features/products/api";
import { useInventoryStream } from "@/features/products/hooks/use-inventory-stream";
import { EventCart } from "@/features/products/ui/EventCart";
import { ProductCard } from "@/features/products/ui/ProductCard";
import { TicketCard } from "@/features/tickets/ui/TicketCard";
import { appQueryClient } from "@/shared/lib/query-client";
import { warmQuery } from "@/shared/lib/route-data";
import { PageNotice } from "@/shared/ui/PageNotice";
import { QueryBoundary } from "@/shared/ui/QueryBoundary";

export const Route = createFileRoute("/events/$slug/store")({
  validateSearch: (search: Record<string, unknown>) => ({
    tab:
      search.tab === "products" ? ("products" as const) : ("tickets" as const),
  }),
  loader: ({ params, context }) =>
    warmQuery(
      appQueryClient,
      storePageQueryOptions(
        params.slug,
        context.session?.isAuthenticated === true,
        appQueryClient,
      ),
    ),
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
  const dataQuery = useQuery(() => storePageQueryOptions(
    params().slug,
    isAuthenticated(),
    queryClient,
  ));
  return (
    <QueryBoundary
      query={dataQuery}
      what="a loja do evento"
      skeleton={<StoreSkeleton />}
      notFound={<PageNotice title="Evento não encontrado" />}
    >
      {(data) => (
        // Keyed, so another event on the same route rebuilds the page instead
        // of leaving the previous event's snapshot on screen.
        <Show
          when={data()}
          keyed
          fallback={<PageNotice title="Evento não encontrado" />}
        >
          {(loaded) => (
            <Show
              when={loaded.edition}
              keyed
              fallback={
                <PageNotice title="A loja não está disponível no momento." />
              }
            >
              {(edition) => (
                <StoreContent
                  event={loaded.event}
                  edition={edition}
                  tickets={loaded.tickets ?? []}
                  products={loaded.products ?? []}
                  initialStock={loaded.stock ?? []}
                  heldTicket={loaded.heldTicket ?? null}
                  initialTab={search().tab}
                />
              )}
            </Show>
          )}
        </Show>
      )}
    </QueryBoundary>
  );
}

/** The store's header, tab bar and card row, at the loaded page's sizes. */
function StoreSkeleton() {
  return (
    <div class="min-h-screen px-4 pb-24 pt-4 md:pt-6">
      <div class="mx-auto max-w-6xl">
        <div class="mx-auto max-w-3xl space-y-3 pt-9 text-center">
          <div class="mx-auto h-4 w-32 animate-pulse rounded bg-muted" />
          <div class="mx-auto h-10 w-3/4 animate-pulse rounded bg-muted" />
          <div class="mx-auto h-5 w-2/3 animate-pulse rounded bg-muted" />
        </div>
        <div class="mt-10 h-12 border-b border-border/60" />
        <div class="mt-8 flex flex-wrap justify-center gap-6 sm:justify-start">
          <For each={[1, 2, 3]}>
            {() => <div class="h-72 w-64 animate-pulse rounded-2xl bg-muted" />}
          </For>
        </div>
      </div>
    </div>
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
    <div class="min-h-screen bg-background px-4 pb-24 pt-4 md:pt-6">
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
        eventSlug={props.event.slug}
      />
    </div>
  );
}
