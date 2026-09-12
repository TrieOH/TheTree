import { createFileRoute, Link } from "@tanstack/solid-router";
import { useQuery } from "@trieoh/front-core-solid";
import { useAuth } from "@trieoh/identityx-sdk-ts-solid";
import {
  For,
  Loading,
  Show,
  createMemo,
  createSignal,
  untrack,
  type Component,
} from "solid-js";
import CheckRaw from "~icons/lucide/check";
import GiftRaw from "~icons/lucide/gift";
import LockRaw from "~icons/lucide/lock";
import MailRaw from "~icons/lucide/mail";
import ShieldCheckRaw from "~icons/lucide/shield-check";
import AlertTriangleRaw from "~icons/lucide/alert-triangle";
import { useCart } from "@/features/products/hooks/use-cart";
import { useCreateCheckoutMutation } from "@/features/purchases/api/mutations";
import {
  checkoutRequest,
  type CheckoutPayment,
} from "@/features/purchases/model/checkout";
import { requireAuth } from "@/features/auths/lib/route-guard";
import { toast } from "@/shared/ui/toast";
import { myTicketQueryOptions } from "@/features/tickets/api";
import {
  checkoutPageQueryOptions,
  type CheckoutPageData,
} from "@/features/purchases/api";
import { OrderSummary } from "@/features/payments/ui/checkout/OrderSummary";
import { PaymentProviderSelector } from "@/features/payments/ui/PaymentProviderSelector";
import { formatMoney } from "@/shared/lib/money";

const LucideCheck = CheckRaw as Component<{ class?: string }>;
const LucideGift = GiftRaw as Component<{ class?: string }>;
const LucideLock = LockRaw as Component<{ class?: string }>;
const LucideMail = MailRaw as Component<{ class?: string }>;
const LucideShieldCheck = ShieldCheckRaw as Component<{ class?: string }>;
const LucideAlertTriangle = AlertTriangleRaw as Component<{ class?: string }>;

export const Route = createFileRoute("/events/$slug/checkout")({
  beforeLoad: requireAuth,
  component: CheckoutPage,
});

function CheckoutPage() {
  const params = Route.useParams();
  const navigate = Route.useNavigate();
  const submitCheckoutMutation = useCreateCheckoutMutation();
  const { auth } = useAuth();
  const [submitting, setSubmitting] = createSignal(false);
  const dataQuery = useQuery<CheckoutPageData | null>(() =>
    checkoutPageQueryOptions(params().slug),
  );
  const data = createMemo(() => dataQuery().data);

  const submitCheckout = async (
    loaded: NonNullable<Awaited<ReturnType<typeof data>>>,
    items: ReturnType<typeof useCart>["items"] extends () => infer T ? T : never,
    payment?: CheckoutPayment,
    gift?: { name: string; email: string },
  ) => {
    const actor = auth.profile();
    if (!actor?.id || !actor.email) {
      toast.error("Complete sua conta antes de comprar.");
      return;
    }
    setSubmitting(true);
    try {
      const result = await submitCheckoutMutation.mutateAsync({
        editionId: loaded.edition.id,
        data: checkoutRequest(items, { id: actor.id, email: actor.email }, payment, gift),
      });
      sessionStorage.setItem(`purchase-ws:${result.purchase_id}`, result.ws_token);
      useCart(loaded.edition.id).clear();
      await navigate({
        to: "/checkouts/$purchaseId",
        params: { purchaseId: result.purchase_id },
      });
    } catch (error) {
      toast.error(
        error instanceof Error
          ? error.message
          : "Não foi possível iniciar o pagamento.",
      );
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Loading fallback={<main class="p-12 text-center">Carregando checkout…</main>}>
      <Show
        when={dataQuery().isSuccess && data()}
        fallback={<main class="p-12 text-center">Evento indisponível.</main>}
      >
        {(loaded) => (
          <CheckoutForm
            editionId={loaded().edition.id}
            eventName={loaded().event.full_name}
            slug={loaded().event.slug}
            submitting={submitting()}
            accountEmail={auth.profile()?.email ?? ""}
            heldTicket={Boolean(loaded().heldTicket)}
            payssagePublicKey={loaded().event.payssage_public_key ?? undefined}
            submit={(items, payment, gift) =>
              submitCheckout(loaded(), items, payment, gift)
            }
          />
        )}
      </Show>
    </Loading>
  );
}

function CheckoutForm(props: {
  editionId: string;
  eventName: string;
  slug: string;
  submitting: boolean;
  accountEmail: string;
  heldTicket: boolean;
  payssagePublicKey?: string;
  submit: (
    items: ReturnType<typeof useCart>["items"] extends () => infer T ? T : never,
    payment?: CheckoutPayment,
    gift?: { name: string; email: string },
  ) => Promise<void>;
}) {
  const cart = useCart(untrack(() => props.editionId));
  const [gift, setGift] = createSignal(false);
  const [giftName, setGiftName] = createSignal("");
  const [giftEmail, setGiftEmail] = createSignal("");

  const hasTicket = () => cart.items().some((item) => item.type === "ticket");
  const totalCents = cart.total;

  // Ticket query to show held-ticket warning
  const ticketQuery = useQuery(() =>
    myTicketQueryOptions(untrack(() => props.editionId)),
  );
  const heldTicket = () => ticketQuery().data ?? null;

  if (cart.items().length === 0) {
    return (
      <main class="mx-auto flex min-h-[70vh] max-w-lg flex-col items-center justify-center gap-4 px-4 text-center">
        <h1 class="text-xl font-semibold">Seu carrinho está vazio</h1>
        <p class="text-sm text-muted-foreground">
          Adicione ao menos um item antes de iniciar o pagamento.
        </p>
        <Link
          to="/events/$slug/store"
          search={{ tab: "products" }}
          params={{ slug: props.slug }}
          class="inline-flex h-9 items-center justify-center bg-primary px-4 text-sm font-medium text-primary-foreground"
        >
          Voltar aos produtos
        </Link>
      </main>
    );
  }

  const handleSubmitPayment = (payment?: CheckoutPayment) => {
    if (hasTicket() && heldTicket() && !gift()) {
      toast.warning("Você já possui um ingresso. Marque a compra como presente para continuar.");
      return;
    }
    if (gift() && !giftName().trim()) {
      toast.error("Informe o nome de quem receberá o presente");
      return;
    }
    if (gift() && !giftEmail().includes("@")) {
      toast.error("Informe um e-mail válido para o presente.");
      return;
    }
    void props.submit(
      cart.items(),
      payment,
      gift()
        ? { name: giftName().trim(), email: giftEmail().trim() }
        : undefined,
    );
  };

  return (
    <main class="mx-auto w-full max-w-6xl px-4 pb-44 pt-6 sm:px-6 md:pb-36 md:pt-12">
      {/* Header + stepper */}
      <div class="mb-6 border-b border-border pb-5 md:mb-8 md:pb-6">
        <div class="mb-3 flex items-center gap-2 text-xs font-medium text-primary md:mb-4">
          <span class="flex size-6 shrink-0 items-center justify-center rounded-full bg-primary text-primary-foreground">
            <LucideCheck class="size-3.5" />
          </span>
          <span class="hidden sm:inline">Pedido conferido</span>
          <span class="h-px w-6 bg-border sm:w-8" />
          <span class="flex size-6 shrink-0 items-center justify-center rounded-full border border-primary">
            2
          </span>
          <span class="text-muted-foreground">
            {totalCents() === 0 ? "Confirmação" : "Pagamento"}
          </span>
        </div>
        <h1 class="text-2xl font-bold tracking-tight sm:text-3xl">
          Finalizar compra
        </h1>
        <p class="mt-1 text-sm text-muted-foreground">
          Confira seu pedido para {props.eventName}.
        </p>
      </div>

      <div class="grid items-start gap-4 lg:grid-cols-[minmax(0,1fr)_25rem] lg:gap-8">
        <div class="space-y-4 lg:space-y-5">
          {/* Order summary */}
          <section class="rounded-md border border-border bg-card p-4 shadow-sm sm:p-5 md:p-6">
            <OrderSummary
              items={cart.items()}
              totalCents={totalCents()}
              title="Seu pedido"
            />
          </section>

          {/* Gift section */}
          <Show when={hasTicket()}>
            <section class="rounded-md border border-border bg-card p-4 shadow-sm sm:p-5 md:p-6">
              <Show when={heldTicket()}>
                <div class="mb-5 flex gap-3 border-l-2 border-amber-500 bg-amber-500/10 p-4 text-sm text-amber-800 dark:text-amber-200">
                  <LucideAlertTriangle class="mt-0.5 size-4 shrink-0" />
                  <p>
                    Você já possui um ingresso.
                    Para comprar outro, marque abaixo que ele é um presente.
                  </p>
                </div>
              </Show>

              <button
                type="button"
                aria-pressed={gift() ? 'true' : 'false'}
                onClick={() => setGift((v) => !v)}
                class="flex w-full items-center gap-3 text-left focus-visible:outline-none sm:gap-4"
              >
                <span class="flex size-10 shrink-0 items-center justify-center rounded-md bg-primary/10 text-primary sm:size-11">
                  <LucideGift class="size-5" />
                </span>
                <span class="flex-1">
                  <span class="block text-sm font-semibold sm:text-base">
                    Este ingresso é um presente?
                  </span>
                  <span class="block text-xs text-muted-foreground sm:text-sm">
                    Enviaremos o ingresso diretamente para quem você escolher.
                  </span>
                </span>
                <span
                  class={`relative h-6 w-11 shrink-0 rounded-full border transition-colors ${gift() ? 'border-primary bg-primary' : 'border-border bg-muted'}`}
                >
                  <span
                    class={`absolute inset-y-0 my-auto size-4 rounded-full shadow-sm ring-1 transition-all ${gift() ? 'translate-x-6 bg-primary-foreground ring-primary-foreground/30' : 'translate-x-1 bg-background ring-border'}`}
                  />
                </span>
              </button>

              <Show when={gift()}>
                <div class="mt-5 grid grid-cols-1 gap-4 border-t border-border pt-5 sm:grid-cols-2">
                  <div class="space-y-2">
                    <label for="gift-name" class="block text-sm font-medium">
                      Nome de quem vai receber
                    </label>
                    <input
                      id="gift-name"
                      required
                      value={giftName()}
                      onInput={(e) => setGiftName(e.currentTarget.value)}
                      placeholder="Nome completo"
                      class="h-10 w-full rounded-md border border-border bg-background px-3 text-sm outline-none focus:border-primary"
                    />
                  </div>
                  <div class="space-y-2">
                    <label for="gift-email" class="block text-sm font-medium">
                      E-mail para entrega
                    </label>
                    <input
                      id="gift-email"
                      type="email"
                      required
                      autocomplete="email"
                      value={giftEmail()}
                      onInput={(e) => setGiftEmail(e.currentTarget.value)}
                      placeholder="presente@exemplo.com"
                      class="h-10 w-full rounded-md border border-border bg-background px-3 text-sm outline-none focus:border-primary"
                    />
                  </div>
                  <p class="flex items-center gap-2 text-xs text-muted-foreground sm:col-span-2">
                    <LucideMail class="size-3.5 shrink-0" /> O presenteado
                    receberá as instruções após a confirmação
                    {totalCents() === 0 ? '.' : ' do pagamento.'}
                  </p>
                </div>
              </Show>
            </section>
          </Show>
        </div>

        {/* Payment sidebar */}
        <section
          id="payment-section"
          class="scroll-mt-20 overflow-hidden rounded-md border border-primary/20 bg-card shadow-md shadow-primary/5 lg:sticky lg:top-6"
        >
          <div class="border-b border-border bg-primary/5 p-5 sm:p-6">
            <div class="mb-2 flex items-center gap-2 text-primary">
              <Show
                when={totalCents() === 0}
                fallback={<LucideLock class="size-4" />}
              >
                <LucideCheck class="size-4" />
              </Show>
              <span class="text-xs font-bold uppercase tracking-wider">
                Finalizar compra
              </span>
            </div>
            <h2 class="text-lg font-bold sm:text-xl">
              {totalCents() === 0
                ? 'Confirme seu pedido'
                : 'Como você prefere pagar?'}
            </h2>
            <p class="mt-1 text-sm text-muted-foreground">
              {totalCents() === 0
                ? 'Nenhum dado de pagamento será necessário.'
                : 'Escolha uma opção para continuar com segurança.'}
            </p>
          </div>

          <div class="p-5 sm:p-6">
            <Show
              when={totalCents() === 0}
              fallback={
                <Show
                  when={props.payssagePublicKey}
                  fallback={
                    <p class="text-sm text-muted-foreground">
                      Este evento ainda não configurou o recebimento de pagamentos.
                    </p>
                  }
                >
                  {(key) => (
                    <PaymentProviderSelector
                      amount={totalCents()}
                      sellerPublicKey={key()}
                      handleSubmit={(payment) => handleSubmitPayment(payment)}
                    />
                  )}
                </Show>
              }
            >
              <div class="flex flex-1 flex-col justify-center gap-4">
                <div class="border-l-2 border-primary bg-primary/5 px-4 py-4">
                  <p class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                    Total do pedido
                  </p>
                  <p class="mt-1 text-2xl font-bold text-primary">Gratuito</p>
                  <p class="mt-2 text-sm text-muted-foreground">
                    Conclua agora sem informar dados de pagamento.
                  </p>
                </div>
                <button
                  type="button"
                  class="h-12 w-full flex items-center justify-between rounded-md bg-primary px-5 text-base font-semibold text-primary-foreground shadow-md shadow-primary/15 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                  disabled={props.submitting}
                  onClick={() => handleSubmitPayment()}
                >
                  <Show when={!props.submitting} fallback="Confirmando…">
                    <span>
                      <span class="sm:hidden">Finalizar grátis</span>
                      <span class="hidden sm:inline">Finalizar pedido gratuito</span>
                    </span>
                    <span class="flex size-6 shrink-0 items-center justify-center rounded-full bg-primary-foreground/15">
                      <LucideCheck class="size-4" />
                    </span>
                  </Show>
                </button>
              </div>
            </Show>
          </div>

          {/* Footer badges */}
          <div class="grid grid-cols-2 gap-3 border-t border-border bg-muted/30 px-5 py-4 text-xs text-muted-foreground sm:px-6">
            <span class="flex items-center gap-2">
              <Show
                when={totalCents() === 0}
                fallback={<LucideShieldCheck class="size-4 shrink-0 text-primary" />}
              >
                <LucideCheck class="size-4 shrink-0 text-primary" />
              </Show>
              {totalCents() === 0 ? 'Sem cobrança' : 'Pagamento protegido'}
            </span>
            <span class="flex items-center gap-2">
              <LucideCheck class="size-4 shrink-0 text-primary" />
              {totalCents() === 0 ? 'Confirmação imediata' : 'Status em tempo real'}
            </span>
          </div>
        </section>
      </div>

      {/* Sticky mobile summary bar */}
      <div class="fixed inset-x-0 bottom-[calc(4rem+env(safe-area-inset-bottom))] z-40 border-t border-border bg-card/95 px-4 py-3 shadow-[0_-4px_16px_rgba(0,0,0,0.06)] backdrop-blur supports-backdrop-filter:bg-card/80 md:hidden">
        <div class="mx-auto flex max-w-5xl items-center justify-between gap-4">
          <div>
            <p class="text-xs text-muted-foreground">Total</p>
            <p class="text-lg font-bold leading-tight">
              {totalCents() === 0 ? 'Gratuito' : formatMoney(totalCents())}
            </p>
          </div>
          <a href="#payment-section">
            <button
              type="button"
              class="h-10 px-6 rounded-md bg-primary text-sm font-semibold text-primary-foreground shadow transition-colors"
            >
              {totalCents() === 0 ? 'Finalizar pedido' : 'Ir para pagamento'}
            </button>
          </a>
        </div>
      </div>
    </main>
  );
}
