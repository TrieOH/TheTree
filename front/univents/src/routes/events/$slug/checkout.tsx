import { useQuery, useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link, redirect } from "@tanstack/react-router";
import { withSpan } from "@trieoh/front-core/tracing/browser";
import { useAuth } from "@trieoh/identityx-sdk-ts/react";
import type {
  CheckoutItem,
  CreateCheckoutRequest,
} from "@trieoh/univents-api/schemas";
import {
  AlertTriangle,
  Check,
  Gift,
  Lock,
  Mail,
  ShieldCheck,
} from "lucide-react";
import { useRef, useState } from "react";
import { toast } from "sonner";
import { requireAuth } from "@/features/auths/lib/route-guard";
import { activeEditionQueryOptions } from "@/features/editions/api";
import { publicEventBySlugQueryOptions } from "@/features/events/api";
import { OrderSummary } from "@/features/payments/ui/checkout/OrderSummary";
import { PaymentProviderSelector } from "@/features/payments/ui/PaymentProviderSelector";
import { useCart } from "@/features/products/hooks/use-cart";
import { profileKeys } from "@/features/profile/api/query-keys";
import { useCreateCheckoutMutation } from "@/features/purchases/api/mutations";
import { myTicketQueryOptions } from "@/features/tickets/api";
import { getErrorMessage } from "@/shared/lib/errors";
import { Button } from "@/shared/ui/shadcn/button";
import { Input } from "@/shared/ui/shadcn/input";
import { Label } from "@/shared/ui/shadcn/label";

function formatBRL(cents: number) {
  return (cents / 100).toLocaleString("pt-BR", {
    style: "currency",
    currency: "BRL",
  });
}

export const Route = createFileRoute("/events/$slug/checkout")({
  beforeLoad: async (args) => {
    requireAuth(args);
    const actorId = args.context.auth?.auth.profile()?.id;
    if (!actorId) return;
    const profile = await args.context.queryClient.ensureQueryData({
      queryKey: profileKeys.detail(actorId),
      queryFn: async () => {
        const response = await args.context.auth?.auth.getActorProfile(actorId);
        return response?.success ? (response.data ?? null) : null;
      },
    });
    if (!profile) {
      throw redirect({
        to: "/profile/setup",
        search: { returnTo: args.location.href },
      });
    }
  },
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
  const hasTicket = items.some((item) => item.type === "ticket");
  const [isGift, setIsGift] = useState(false);
  const currentEditionId = activeEdition?.id ?? edition.id;
  const {
    data: heldTicket = null,
    isPending: isCheckingTicket,
    isError: ticketCheckFailed,
  } = useQuery(myTicketQueryOptions(currentEditionId, true));
  const isUpgrade = items.some(
    (item) => item.type === "ticket" && item.is_upgrade,
  );
  const [giftName, setGiftName] = useState("");
  const [giftEmail, setGiftEmail] = useState("");
  const giftEmailRef = useRef<HTMLInputElement>(null);

  const submitPayment = async (payment?: {
    card_token?: string;
    payment_method_id: string;
    payment_method_type: string;
    issuer_id?: string;
    installments?: number;
    payer_email: string;
    identification_type: string;
    identification_number: string;
  }) => {
    if (hasTicket && isCheckingTicket) {
      toast.info("Aguarde enquanto verificamos seu ingresso atual");
      return;
    }
    if (hasTicket && ticketCheckFailed) {
      toast.error(
        "Não foi possível verificar seu ingresso atual. Tente novamente.",
      );
      return;
    }
    if (heldTicket && !isGift && hasTicket) {
      toast.warning(
        isUpgrade
          ? "Upgrade de ingresso estará disponível em breve. Se este ingresso for para outra pessoa, marque-o como presente."
          : "Você já possui um ingresso. Marque a compra como presente para continuar.",
      );
      return;
    }
    if (isGift && !giftName.trim()) {
      toast.error("Informe o nome de quem receberá o presente");
      return;
    }
    if (isGift && !giftEmailRef.current?.checkValidity()) {
      giftEmailRef.current?.reportValidity();
      return;
    }

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
            user_id: isGift ? undefined : buyer.id,
            email: isGift ? giftEmail.trim() : buyer.email,
            name: isGift ? giftName.trim() : buyer.email,
          },
        });
      }
    }

    const data: CreateCheckoutRequest =
      totalCents === 0
        ? { items: checkoutItems }
        : {
            payment_method:
              !payment || payment.payment_method_id === "pix"
                ? "pix"
                : "credit_card",
            card_token: payment?.card_token,
            payment_method_id:
              !payment || payment.payment_method_id === "pix"
                ? undefined
                : payment.payment_method_id,
            issuer_id: payment?.issuer_id,
            installments: payment?.installments,
            payer: {
              email: payment?.payer_email ?? profile.email,
              identification_type: payment?.identification_type ?? "",
              identification_number: payment?.identification_number ?? "",
            },
            items: checkoutItems,
          };

    try {
      const result = await withSpan("action:checkout", () =>
        checkout.mutateAsync({ editionId: edition.id, data }),
      );
      sessionStorage.setItem(
        `purchase-ws:${result.purchase_id}`,
        result.ws_token,
      );
      clearCart();
      await navigate({
        to: "/checkouts/$purchaseId",
        params: { purchaseId: result.purchase_id },
        search: { state: "pending" },
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
          to="/events/$slug/store"
          search={{ tab: "products" }}
          params={{ slug: event.slug }}
          className="inline-flex h-9 items-center justify-center bg-primary px-4 text-sm font-medium text-primary-foreground"
        >
          Voltar aos produtos
        </Link>
      </main>
    );
  }

  return (
    <main className="mx-auto w-full max-w-6xl px-4 pb-44 pt-6 sm:px-6 md:pb-36 md:pt-12">
      {/* Header + stepper */}
      <div className="mb-6 border-b border-border pb-5 md:mb-8 md:pb-6">
        <div className="mb-3 flex items-center gap-2 text-xs font-medium text-primary md:mb-4">
          <span className="flex size-6 shrink-0 items-center justify-center rounded-full bg-primary text-primary-foreground">
            <Check className="size-3.5" />
          </span>
          <span className="hidden sm:inline">Pedido conferido</span>
          <span className="h-px w-6 bg-border sm:w-8" />
          <span className="flex size-6 shrink-0 items-center justify-center rounded-full border border-primary">
            2
          </span>
          <span className="text-muted-foreground">
            {totalCents === 0 ? "Confirmação" : "Pagamento"}
          </span>
        </div>
        <h1 className="text-2xl font-bold tracking-tight sm:text-3xl">
          Finalizar compra
        </h1>
        <p className="mt-1 text-sm text-muted-foreground">
          Confira seu pedido para {event.full_name}.
        </p>
      </div>

      <div className="grid items-start gap-4 lg:grid-cols-[minmax(0,1fr)_25rem] lg:gap-8">
        <div className="space-y-4 lg:space-y-5">
          <section className="rounded-md border border-border bg-card p-4 shadow-sm sm:p-5 md:p-6">
            <OrderSummary
              items={items}
              totalCents={totalCents}
              title="Seu pedido"
            />
          </section>

          {hasTicket && (
            <section className="rounded-md border border-border bg-card p-4 shadow-sm sm:p-5 md:p-6">
              {heldTicket && (
                <div className="mb-5 flex gap-3 border-l-2 border-amber-500 bg-amber-500/10 p-4 text-sm text-amber-800 dark:text-amber-200">
                  <AlertTriangle className="mt-0.5 size-4 shrink-0" />
                  <p>
                    Você já possui o ingresso {heldTicket.ticket_type.name}.
                    Para comprar outro, marque abaixo que ele é um presente.
                  </p>
                </div>
              )}
              <button
                type="button"
                aria-pressed={isGift}
                onClick={() => setIsGift((value) => !value)}
                className="flex w-full items-center gap-3 text-left focus-visible:outline-none sm:gap-4"
              >
                <span className="flex size-10 shrink-0 items-center justify-center rounded-md bg-primary/10 text-primary sm:size-11">
                  <Gift className="size-5" />
                </span>
                <span className="flex-1">
                  <span className="block text-sm font-semibold sm:text-base">
                    Este ingresso é um presente?
                  </span>
                  <span className="block text-xs text-muted-foreground sm:text-sm">
                    Enviaremos o ingresso diretamente para quem você escolher.
                  </span>
                </span>
                <span
                  className={`relative h-6 w-11 shrink-0 rounded-full border transition-colors ${isGift ? "border-primary bg-primary" : "border-border bg-muted"}`}
                >
                  <span
                    className={`absolute inset-y-0 my-auto size-4 rounded-full shadow-sm ring-1 transition-all ${isGift ? "translate-x-6 bg-primary-foreground ring-primary-foreground/30" : "translate-x-1 bg-background ring-border"}`}
                  />
                </span>
              </button>

              {isGift && (
                <div className="mt-5 grid grid-cols-1 gap-4 border-t border-border pt-5 sm:grid-cols-2">
                  <div className="space-y-2">
                    <Label htmlFor="gift-name">Nome de quem vai receber</Label>
                    <Input
                      id="gift-name"
                      required
                      value={giftName}
                      onChange={(event) => setGiftName(event.target.value)}
                      placeholder="Nome completo"
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="gift-email">E-mail para entrega</Label>
                    <Input
                      id="gift-email"
                      ref={giftEmailRef}
                      type="email"
                      required
                      autoComplete="email"
                      value={giftEmail}
                      onChange={(event) => setGiftEmail(event.target.value)}
                      placeholder="presente@exemplo.com"
                    />
                  </div>
                  <p className="flex items-center gap-2 text-xs text-muted-foreground sm:col-span-2">
                    <Mail className="size-3.5 shrink-0" /> O presenteado
                    receberá as instruções após a confirmação
                    {totalCents === 0 ? "." : " do pagamento."}
                  </p>
                </div>
              )}
            </section>
          )}
        </div>

        <section
          id="payment-section"
          className="scroll-mt-20 overflow-hidden rounded-md border border-primary/20 bg-card shadow-md shadow-primary/5 lg:sticky lg:top-6"
        >
          <div className="border-b border-border bg-primary/5 p-5 sm:p-6">
            <div className="mb-2 flex items-center gap-2 text-primary">
              {totalCents === 0 ? (
                <Check className="size-4" />
              ) : (
                <Lock className="size-4" />
              )}
              <span className="text-xs font-bold uppercase tracking-wider">
                Finalizar compra
              </span>
            </div>
            <h2 className="text-lg font-bold sm:text-xl">
              {totalCents === 0
                ? "Confirme seu pedido"
                : "Como você prefere pagar?"}
            </h2>
            <p className="mt-1 text-sm text-muted-foreground">
              {totalCents === 0
                ? "Nenhum dado de pagamento será necessário."
                : "Escolha uma opção para continuar com segurança."}
            </p>
          </div>

          <div className="p-5 sm:p-6">
            {totalCents === 0 ? (
              <div className="flex flex-1 flex-col justify-center gap-4">
                <div className="border-l-2 border-primary bg-primary/5 px-4 py-4">
                  <p className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                    Total do pedido
                  </p>
                  <p className="mt-1 text-2xl font-bold text-primary">
                    Gratuito
                  </p>
                  <p className="mt-2 text-sm text-muted-foreground">
                    Conclua agora sem informar dados de pagamento.
                  </p>
                </div>
                <Button
                  size="lg"
                  className="h-12 w-full justify-between px-5 text-base shadow-md shadow-primary/15"
                  disabled={checkout.isPending}
                  onClick={() => void submitPayment()}
                >
                  {checkout.isPending ? (
                    "Confirmando…"
                  ) : (
                    <span>
                      <span className="sm:hidden">Finalizar grátis</span>
                      <span className="hidden sm:inline">
                        Finalizar pedido gratuito
                      </span>
                    </span>
                  )}
                  {!checkout.isPending && (
                    <span className="flex size-6 shrink-0 items-center justify-center rounded-full bg-primary-foreground/15">
                      <Check className="size-4" />
                    </span>
                  )}
                </Button>
              </div>
            ) : event.payssage_public_key ? (
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
          </div>

          <div className="grid grid-cols-2 gap-3 border-t border-border bg-muted/30 px-5 py-4 text-xs text-muted-foreground sm:px-6">
            <span className="flex items-center gap-2">
              {totalCents === 0 ? (
                <Check className="size-4 shrink-0 text-primary" />
              ) : (
                <ShieldCheck className="size-4 shrink-0 text-primary" />
              )}
              {totalCents === 0 ? "Sem cobrança" : "Pagamento protegido"}
            </span>
            <span className="flex items-center gap-2">
              <Check className="size-4 shrink-0 text-primary" />
              {totalCents === 0
                ? "Confirmação imediata"
                : "Status em tempo real"}
            </span>
          </div>
        </section>
      </div>

      {/* Sticky mobile summary bar — keeps total + CTA reachable without
          scrolling all the way down, and jumps straight to payment. */}
      <div className="fixed inset-x-0 bottom-[calc(4rem+env(safe-area-inset-bottom))] z-40 border-t border-border bg-card/95 px-4 py-3 shadow-[0_-4px_16px_rgba(0,0,0,0.06)] backdrop-blur supports-backdrop-filter:bg-card/80 md:hidden">
        <div className="mx-auto flex max-w-5xl items-center justify-between gap-4">
          <div>
            <p className="text-xs text-muted-foreground">Total</p>
            <p className="text-lg font-bold leading-tight">
              {totalCents === 0 ? "Gratuito" : formatBRL(totalCents)}
            </p>
          </div>
          <a href="#payment-section">
            <Button size="lg" className="px-6">
              {totalCents === 0 ? "Finalizar pedido" : "Ir para pagamento"}
            </Button>
          </a>
        </div>
      </div>
    </main>
  );
}
