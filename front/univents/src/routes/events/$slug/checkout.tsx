import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link, redirect } from "@tanstack/react-router";
import { useAuth } from "@trieoh/identityx-sdk-ts/react";
import type {
  CheckoutItem,
  CreateCheckoutRequest,
} from "@trieoh/univents-api/schemas";
import { Check } from "lucide-react";
import { toast } from "sonner";
import { requireAuth } from "@/features/auths/lib/route-guard";
import { activeEditionQueryOptions } from "@/features/editions/api";
import { publicEventBySlugQueryOptions } from "@/features/events/api";
import { OrderSummary } from "@/features/payments/ui/checkout/OrderSummary";
import { PaymentProviderSelector } from "@/features/payments/ui/PaymentProviderSelector";
import { useCart } from "@/features/products/hooks/use-cart";
import { useCreateCheckoutMutation } from "@/features/purchases/api";
import { getErrorMessage } from "@/shared/lib/errors";

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
  const navigate = Route.useNavigate();
  const { event, edition } = Route.useLoaderData();
  const { auth } = useAuth();
  const { data: activeEdition } = useSuspenseQuery(
    activeEditionQueryOptions(event.id),
  );
  const { items, totalCents, clearCart } = useCart(
    activeEdition?.id ?? edition.id,
  );
  const checkout = useCreateCheckoutMutation();

  const submitPayment = async (payment: {
    card_token?: string;
    payment_method_id: string;
    payment_method_type: string;
    installments?: number;
    payer_email: string;
    identification_type: string;
    identification_number: string;
  }) => {
    const profile = auth.profile();
    if (!profile?.id || !profile.email) {
      toast.error("Complete seu perfil antes de comprar");
      return;
    }

    const buyer = { id: profile.id, email: profile.email };
    const checkoutItems: CheckoutItem[] = [];
    for (const item of items) {
      const itemType =
        item.type === "activity" ? "program_occurrence" : item.type;
      if (itemType !== "ticket") {
        checkoutItems.push({
          item_type: itemType,
          item_id: item.id,
          quantity: item.quantity,
        });
        continue;
      }
      for (let quantity = 0; quantity < item.quantity; quantity++) {
        checkoutItems.push({
          item_type: "ticket",
          item_id: item.id,
          quantity: 1,
          attendee: {
            user_id: buyer.id,
            email: buyer.email,
            name: buyer.email,
          },
        });
      }
    }

    const data: CreateCheckoutRequest = {
      payment_method:
        payment.payment_method_id === "pix" ? "pix" : "credit_card",
      card_token: payment.card_token,
      payment_method_id:
        payment.payment_method_id === "pix"
          ? undefined
          : payment.payment_method_id,
      installments: payment.installments,
      payer: {
        email: payment.payer_email,
        identification_type: payment.identification_type,
        identification_number: payment.identification_number,
      },
      items: checkoutItems,
    };

    try {
      const result = await checkout.mutateAsync({
        editionId: edition.id,
        data,
      });
      sessionStorage.setItem(
        `purchase-ws:${result.purchase_id}`,
        result.ws_token,
      );
      clearCart();
      await navigate({
        to: "/checkouts/$purchaseId",
        params: { purchaseId: result.purchase_id },
      });
    } catch (error) {
      toast.error(
        getErrorMessage(error, "Não foi possível iniciar o pagamento"),
      );
    }
  };

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
          {event.payssage_public_key ? (
            <PaymentProviderSelector
              amount={totalCents}
              sellerPublicKey={event.payssage_public_key}
              handleSubmit={(payment) => void submitPayment(payment)}
            />
          ) : (
            <p className="text-sm text-muted-foreground">
              Este evento ainda não configurou o recebimento de pagamentos.
            </p>
          )}
        </section>
      </div>
    </main>
  );
}
