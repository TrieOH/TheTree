import { useQueries, useQuery, useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { useAuth } from "@trieoh/identityx-sdk-ts/react";
import { ArrowLeft } from "lucide-react";
import { useState } from "react";
import { activeEditionQueryOptions } from "@/features/editions/api";
import { publicEventBySlugQueryOptions } from "@/features/events/api";
import {
  productsByEditionQueryOptions,
  productVariantsQueryOptions,
} from "@/features/products/api";
import { EventCart } from "@/features/products/ui/EventCart";
import { ProductCardCompact } from "@/features/products/ui/ProductCardCompact";
import {
  allTicketsQueryOptions,
  myTicketQueryOptions,
} from "@/features/tickets/api";
import { TicketCard } from "@/features/tickets/ui/TicketCard";

export const Route = createFileRoute("/events/$slug/store")({
  validateSearch: (search: Record<string, unknown>) => ({
    tab: search.tab === "products" ? "products" : "tickets",
  }),
  loader: async ({ context, params }) => {
    const event = await context.queryClient.ensureQueryData(
      publicEventBySlugQueryOptions(params.slug),
    );
    if (!event) return null;
    return event;
  },
  component: StorePage,
});

function StorePage() {
  const event = Route.useLoaderData();
  const { tab } = Route.useSearch();
  const navigate = Route.useNavigate();
  const { isAuthenticated } = useAuth();
  const [activeTab, setActiveTab] = useState(tab);

  if (!event)
    return <main className="p-12 text-center">Evento não encontrado.</main>;

  const { data: edition } = useSuspenseQuery(
    activeEditionQueryOptions(event.id),
  );
  const { data: tickets = [] } = useSuspenseQuery(
    allTicketsQueryOptions(edition?.id ?? ""),
  );
  const { data: products = [] } = useSuspenseQuery(
    productsByEditionQueryOptions(edition?.id ?? ""),
  );
  const variants = useQueries({
    queries: products.map((product) => productVariantsQueryOptions(product.id)),
  });
  const { data: heldTicket = null } = useQuery(
    myTicketQueryOptions(edition?.id ?? "", isAuthenticated),
  );

  if (!edition) {
    return (
      <main className="mx-auto max-w-3xl px-4 py-16 text-center">
        A loja não está disponível no momento.
      </main>
    );
  }

  const setTab = (next: "tickets" | "products") => {
    setActiveTab(next);
    void navigate({ search: { tab: next } });
  };

  return (
    <main className="min-h-screen bg-background px-4 pt-4 pb-24 md:pt-6">
      <div className="mx-auto max-w-6xl">
        <header className="mx-auto max-w-3xl">
          <Link
            to="/events/$slug"
            params={{ slug: event.slug }}
            className="inline-flex items-center gap-2 rounded-lg px-2 py-2 text-sm font-medium text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
          >
            <ArrowLeft className="size-4" /> Evento
          </Link>
          <div className="mt-5 text-center">
            <p className="text-xs font-semibold uppercase tracking-[0.2em] text-primary">
              Loja do evento
            </p>
            <h1 className="mt-3 text-3xl font-bold tracking-tight sm:text-4xl">
              Ingressos e produtos
            </h1>
            <p className="mx-auto mt-3 max-w-xl text-muted-foreground">
              Encontre ingressos e produtos para {event.full_name}.
            </p>
          </div>
        </header>
        <div className="mt-10 flex gap-2 border-b border-border/60">
          {(["tickets", "products"] as const).map((item) => (
            <button
              key={item}
              type="button"
              onClick={() => setTab(item)}
              className={`border-b-2 px-4 py-3 text-sm font-medium ${activeTab === item ? "border-primary text-primary" : "border-transparent text-muted-foreground"}`}
            >
              {item === "tickets" ? "Ingressos" : "Produtos"}
            </button>
          ))}
        </div>
        {activeTab === "tickets" ? (
          tickets.length ? (
            <div className="mt-8 flex flex-wrap gap-6">
              {tickets.map((ticket, index) => (
                <TicketCard
                  key={ticket.id}
                  ticket={ticket}
                  editionId={edition.id}
                  heldTicket={heldTicket}
                  isFeatured={index === 0}
                />
              ))}
            </div>
          ) : (
            <p className="mt-8 text-muted-foreground">
              Nenhum ingresso disponível.
            </p>
          )
        ) : products.some(
            (_, index) => (variants[index]?.data ?? []).length > 0,
          ) ? (
          <div className="mt-8 flex flex-wrap gap-6">
            {products.map((product, index) => {
              const productVariants = variants[index]?.data ?? [];
              return productVariants.length ? (
                <ProductCardCompact
                  key={product.id}
                  product={product}
                  variants={productVariants}
                  editionId={edition.id}
                />
              ) : null;
            })}
          </div>
        ) : (
          <div className="mt-8 rounded-xl border border-dashed p-10 text-center">
            <h2 className="text-lg font-semibold">Nenhum produto disponível</h2>
            <p className="mt-2 text-sm text-muted-foreground">
              Ainda não há produtos publicados para esta edição.
            </p>
            <Link
              to="/events/$slug"
              params={{ slug: event.slug }}
              className="mt-6 inline-flex text-sm font-medium text-primary hover:underline"
            >
              Voltar ao evento
            </Link>
          </div>
        )}
      </div>
      <EventCart
        editionId={edition.id}
        onCheckout={() =>
          navigate({
            to: "/events/$slug/checkout",
            params: { slug: event.slug },
          })
        }
        onExplore={() => setTab("products")}
      />
    </main>
  );
}
