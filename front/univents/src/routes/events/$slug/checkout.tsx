import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link, redirect } from "@tanstack/react-router";
import { ArrowLeft, LockKeyhole } from "lucide-react";
import { requireAuth } from "@/features/auths/lib/route-guard";
import { activeEditionQueryOptions } from "@/features/editions/api";
import { publicEventBySlugQueryOptions } from "@/features/events/api";
import { OrderSummary } from "@/features/payments/ui/checkout/OrderSummary";
import { useCart } from "@/features/products/hooks/use-cart";
import { Button } from "@/shared/ui/shadcn/button";

export const Route = createFileRoute("/events/$slug/checkout")({
  beforeLoad: requireAuth,
  loader: async ({ context, params }) => {
    const event = await context.queryClient.ensureQueryData(
      publicEventBySlugQueryOptions(params.slug),
    );
    if (!event) throw redirect({ to: "/events" });

    const edition = await context.queryClient.ensureQueryData(
      activeEditionQueryOptions(event.id),
    );
    if (!edition) {
      throw redirect({ to: "/events/$slug", params: { slug: params.slug } });
    }

    return { event, edition };
  },
  component: CheckoutPage,
});

function CheckoutPage() {
  const { event, edition } = Route.useLoaderData();
  const { data: activeEdition } = useSuspenseQuery(
    activeEditionQueryOptions(event.id),
  );
  const { items, totalCents } = useCart(activeEdition?.id ?? edition.id);

  if (items.length === 0) {
    return (
      <main className="mx-auto flex min-h-[70vh] max-w-lg flex-col items-center justify-center gap-4 px-4 text-center">
        <h1 className="text-xl font-semibold">Seu carrinho está vazio</h1>
        <p className="text-sm text-muted-foreground">
          Adicione ao menos um item antes de iniciar o pagamento.
        </p>
        <Link
          to="/events/$slug/products"
          params={{ slug: event.slug }}
          className="inline-flex h-9 items-center justify-center bg-primary px-4 text-sm font-medium text-primary-foreground"
        >
          Voltar aos produtos
        </Link>
      </main>
    );
  }

  return (
    <main className="mx-auto w-full max-w-4xl px-4 py-8 md:py-12">
      <Link
        to="/events/$slug/products"
        params={{ slug: event.slug }}
        className="mb-6 inline-flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground"
      >
        <ArrowLeft className="size-4" />
        Voltar ao carrinho
      </Link>

      <div className="mb-8">
        <h1 className="text-2xl font-bold">Checkout</h1>
        <p className="mt-1 text-sm text-muted-foreground">
          Confira seu pedido para {event.full_name}.
        </p>
      </div>

      <div className="grid gap-8 lg:grid-cols-2">
        <OrderSummary items={items} totalCents={totalCents} />

        <section className="flex flex-col justify-center border border-border bg-card p-6 text-center">
          <LockKeyhole className="mx-auto mb-3 size-8 text-primary" />
          <h2 className="font-semibold">Pagamento seguro</h2>
          <p className="mt-2 text-sm text-muted-foreground">
            A seleção de Pix ou cartão será liberada aqui quando o endpoint de
            checkout dos próximos splits estiver disponível.
          </p>
          <Button className="mt-6" disabled>
            Continuar para pagamento
          </Button>
        </section>
      </div>
    </main>
  );
}
