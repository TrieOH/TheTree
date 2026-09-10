import { createFileRoute, Link } from "@tanstack/solid-router";
import { useQueryClient } from "@trieoh/front-core/solid";
import { useAuth } from "@trieoh/identityx-sdk-ts-solid";
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
import { useCart } from "@/features/products/hooks/use-cart";
import { createCheckout } from "@/features/purchases/api";
import {
  checkoutRequest,
  type CheckoutPayment,
} from "@/features/purchases/model/checkout";
import { requireAuth } from "@/features/auths/lib/route-guard";
import { toast } from "@/shared/ui/toast";
import { myTicketQueryOptions } from "@/features/tickets/api";

export const Route = createFileRoute("/events/$slug/checkout")({
  beforeLoad: requireAuth,
  component: CheckoutPage,
});

function CheckoutPage() {
  const params = Route.useParams();
  const navigate = Route.useNavigate();
  const queryClient = useQueryClient();
  const { auth } = useAuth();
  const [submitting, setSubmitting] = createSignal(false);
  const data = createMemo(async () => {
    const event = await queryClient.fetchQuery(
      publicEventBySlugQueryOptions(params().slug),
    );
    if (!event) return null;
    const editions = await queryClient.fetchQuery(
      allPublicEditionsQueryOptions(event.id),
    );
    const now = Date.now();
    const sorted = [...editions].sort((a, b) =>
      a.starts_at.localeCompare(b.starts_at),
    );
    const edition =
      sorted.find(
        (candidate) =>
          new Date(candidate.starts_at).getTime() <= now &&
          new Date(candidate.ends_at).getTime() >= now,
      ) ??
      sorted.find((candidate) => new Date(candidate.starts_at).getTime() > now) ??
      sorted.at(-1);
    if (!edition) return null;
    const heldTicket = await queryClient.fetchQuery(myTicketQueryOptions(edition.id));
    return { event, edition, heldTicket };
  });
  const submitCheckout = async (
    loaded: NonNullable<Awaited<ReturnType<typeof data>>>,
    items: ReturnType<typeof useCart>["items"] extends () => infer T
      ? T
      : never,
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
      const result = await createCheckout(
        loaded.edition.id,
        checkoutRequest(
          items,
          { id: actor.id, email: actor.email },
          payment,
          gift,
        ),
      );
      sessionStorage.setItem(
        `purchase-ws:${result.purchase_id}`,
        result.ws_token,
      );
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
    <Loading
      fallback={<main class="p-12 text-center">Carregando checkout…</main>}
    >
      <Show
        when={data()}
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
          paymentsConfigured={Boolean(loaded().event.payssage_public_key)}
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
  paymentsConfigured: boolean;
  submit: (
    items: ReturnType<typeof useCart>["items"] extends () => infer T
      ? T
      : never,
    payment?: CheckoutPayment,
    gift?: { name: string; email: string },
  ) => Promise<void>;
}) {
  const cart = useCart(untrack(() => props.editionId));
  const [gift, setGift] = createSignal(false);
  const [giftName, setGiftName] = createSignal("");
  const [giftEmail, setGiftEmail] = createSignal("");
  const [payerEmail, setPayerEmail] = createSignal(
    untrack(() => props.accountEmail),
  );
  const [documentType, setDocumentType] = createSignal("CPF");
  const [documentNumber, setDocumentNumber] = createSignal("");
  const hasTicket = () => cart.items().some((item) => item.type === "ticket");
  const money = (cents: number) =>
    new Intl.NumberFormat("pt-BR", {
      style: "currency",
      currency: "BRL",
    }).format(cents / 100);
  const finish = () => {
    if (hasTicket() && props.heldTicket && !gift())
      return toast.error("Você já possui um ingresso. Marque a compra como presente para continuar.");
    if (gift() && (!giftName().trim() || !giftEmail().includes("@")))
      return toast.error("Informe nome e e-mail válidos para o presente.");
    if (
      cart.total() > 0 &&
      (!payerEmail().includes("@") ||
        documentNumber().replace(/\D/g, "").length < 11)
    )
      return toast.error("Informe e-mail e CPF/CNPJ válidos para gerar o Pix.");
    const payment =
      cart.total() > 0
        ? {
            method: "pix" as const,
            email: payerEmail().trim(),
            identificationType: documentType(),
            identificationNumber: documentNumber().replace(/\D/g, ""),
          }
        : undefined;
    void props.submit(
      cart.items(),
      payment,
      gift()
        ? { name: giftName().trim(), email: giftEmail().trim() }
        : undefined,
    );
  };
  return (
    <main class="mx-auto max-w-6xl px-4 pb-40 pt-8">
      <Link
        to="/events/$slug/store"
        params={{ slug: props.slug }}
        search={{ tab: "products" }}
        class="text-sm text-muted-foreground"
      >
        ← Voltar à loja
      </Link>
      <header class="mt-6 border-b border-border pb-6">
        <p class="text-xs font-semibold text-primary">
          ✓ Pedido conferido — 2 Pagamento
        </p>
        <h1 class="mt-3 text-3xl font-bold">Finalizar compra</h1>
        <p class="mt-1 text-sm text-muted-foreground">
          Confira seu pedido para {props.eventName}.
        </p>
      </header>
      <Show
        when={cart.items().length}
        fallback={
          <p class="mt-10 rounded-md border border-border bg-card p-6">
            Seu carrinho está vazio.
          </p>
        }
      >
        <div class="mt-8 grid items-start gap-6 lg:grid-cols-[minmax(0,1fr)_360px]">
          <div class="space-y-5">
            <section class="rounded-md border border-border bg-card p-5">
              <h2 class="mb-4 font-semibold">Seu pedido</h2>
              <For each={cart.items()}>
                {(item) => (
                  <div class="flex justify-between gap-4 border-b border-border py-4 first:pt-0 last:border-0">
                    <div>
                      <strong>{item.name}</strong>
                      <p class="text-sm text-muted-foreground">
                        {item.quantity} unidade(s)
                      </p>
                    </div>
                    <span>{money(item.price_cents * item.quantity)}</span>
                  </div>
                )}
              </For>
            </section>
              <Show when={hasTicket()}>
                <section class="rounded-md border border-border bg-card p-5">
                  <Show when={props.heldTicket && !gift()}>
                    <p class="mb-4 border-l-2 border-amber-500 bg-amber-500/10 p-3 text-sm text-amber-700">Você já possui um ingresso para esta edição.</p>
                  </Show>
                <button
                  type="button"
                  aria-pressed={gift() ? "true" : "false"}
                  class="flex w-full items-center justify-between text-left"
                  onClick={() => setGift((value) => !value)}
                >
                  <span>
                    <strong>Este ingresso é um presente?</strong>
                    <small class="mt-1 block text-muted-foreground">
                      Enviaremos diretamente para quem você escolher.
                    </small>
                  </span>
                  <span
                    class={`relative h-6 w-11 rounded-full ${gift() ? "bg-primary" : "bg-muted"}`}
                  >
                    <i
                      class={`absolute top-1 size-4 rounded-full bg-background transition-transform ${gift() ? "translate-x-6" : "translate-x-1"}`}
                    />
                  </span>
                </button>
                <Show when={gift()}>
                  <div class="mt-5 grid gap-4 border-t border-border pt-5 sm:grid-cols-2">
                    <Field
                      label="Nome de quem vai receber"
                      value={giftName()}
                      onInput={setGiftName}
                    />
                    <Field
                      label="E-mail para entrega"
                      type="email"
                      value={giftEmail()}
                      onInput={setGiftEmail}
                    />
                  </div>
                </Show>
              </section>
            </Show>
          </div>
          <aside
            id="payment-section"
            class="rounded-md border border-primary/20 bg-card shadow-md lg:sticky lg:top-6"
          >
            <div class="border-b border-border bg-primary/5 p-5">
              <p class="text-xs font-bold uppercase text-primary">
                Finalizar compra
              </p>
              <h2 class="mt-2 text-xl font-bold">
                {cart.total() ? "Pagar com Pix" : "Confirme seu pedido"}
              </h2>
            </div>
            <div class="space-y-4 p-5">
              <div>
                <p class="text-sm text-muted-foreground">Total</p>
                <strong class="text-3xl text-primary">
                  {cart.total() ? money(cart.total()) : "Gratuito"}
                </strong>
              </div>
              <Show when={cart.total() > 0}>
                <Field
                  label="E-mail do pagador"
                  type="email"
                  value={payerEmail()}
                  onInput={setPayerEmail}
                />
                <div class="grid grid-cols-[110px_1fr] gap-3">
                  <label class="space-y-2 text-sm font-medium">
                    Documento
                    <select
                      value={documentType()}
                      onChange={(event) =>
                        setDocumentType(event.currentTarget.value)
                      }
                      class="h-10 w-full rounded-md border border-border bg-background px-3"
                    >
                      <option>CPF</option>
                      <option>CNPJ</option>
                    </select>
                  </label>
                  <Field
                    label={documentType()}
                    value={documentNumber()}
                    onInput={setDocumentNumber}
                  />
                </div>
              </Show>
                <button
                  type="button"
                  disabled={props.submitting || (cart.total() > 0 && !props.paymentsConfigured)}
                class="h-12 w-full rounded-md bg-primary font-semibold text-primary-foreground disabled:opacity-50"
                onClick={finish}
              >
                {props.submitting
                  ? "Processando…"
                  : cart.total()
                    ? "Gerar QR Code Pix"
                    : "Finalizar pedido gratuito"}
                </button>
                <Show when={cart.total() > 0 && !props.paymentsConfigured}>
                  <p class="text-sm text-destructive">Este evento ainda não configurou o recebimento de pagamentos.</p>
                </Show>
              <p class="text-center text-xs text-muted-foreground">
                Pagamento protegido · status em tempo real
              </p>
            </div>
          </aside>
        </div>
      </Show>
    </main>
  );
}

function Field(props: {
  label: string;
  value: string;
  onInput: (value: string) => void;
  type?: "text" | "email";
}) {
  return (
    <label class="space-y-2 text-sm font-medium">
      {props.label}
      <input
        type={props.type ?? "text"}
        value={props.value}
        onInput={(event) => props.onInput(event.currentTarget.value)}
        class="h-10 w-full rounded-md border border-border bg-background px-3 outline-none focus:border-primary"
      />
    </label>
  );
}
