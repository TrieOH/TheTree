import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link, redirect } from "@tanstack/react-router";
import { Check, CreditCard, LockKeyhole } from "lucide-react";
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
    <main className="mx-auto w-full max-w-5xl px-4 py-8 md:py-12">
      <div className="mb-8 border-b border-border pb-6">
        <div className="mb-4 flex items-center gap-2 text-xs font-medium text-primary">
          <span className="flex size-6 items-center justify-center rounded-full bg-primary text-primary-foreground">
            <Check className="size-3.5" />
          </span>
          Pedido conferido
          <span className="h-px w-8 bg-border" />
          <span className="flex size-6 items-center justify-center rounded-full border border-primary">
            2
          </span>
          <span className="text-muted-foreground">Pagamento</span>
        </div>
        <h1 className="text-3xl font-bold tracking-tight">Finalizar compra</h1>
        <p className="mt-1 text-sm text-muted-foreground">
          Confira seu pedido para {event.full_name}.
        </p>
      </div>

      <div className="grid gap-8 lg:grid-cols-2">
        <section className="rounded-xl border border-border bg-card p-5 md:p-6">
          <OrderSummary
            items={items}
            totalCents={totalCents}
            title="Seu pedido"
          />
        </section>

        <section className="flex flex-col rounded-xl border border-primary/20 bg-primary/3 p-6">
          <div className="mb-6 flex items-start gap-3">
            <div className="flex size-10 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
              <CreditCard className="size-5" />
            </div>
            <div>
              <h2 className="font-semibold">Forma de pagamento</h2>
              <p className="mt-1 text-xs text-muted-foreground">
                Pix e cartão de crédito
              </p>
            </div>
            <LockKeyhole className="ml-auto size-4 text-muted-foreground" />
          </div>
          <p className="mt-2 text-sm text-muted-foreground">
            O pagamento estará disponível assim que a integração financeira
            deste evento for liberada.
          </p>
          <Button className="mt-6 w-full" disabled>
            Continuar para o pagamento
          </Button>
          <p className="mt-3 text-center text-[11px] text-muted-foreground">
            Seus dados serão processados com segurança.
          </p>
        </section>
      </div>
    </main>
  );
}
