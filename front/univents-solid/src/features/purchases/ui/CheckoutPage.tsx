import {
  createEffect,
  createMemo,
  createSignal,
  For,
  Loading,
  Show,
  type Component,
} from "solid-js";
import { Link } from "@tanstack/solid-router";
import type { Purchase } from "@trieoh/univents-api/schemas";
import { useQueryClient } from "@trieoh/front-core-solid";
import QRCode from "qrcode";

import { formatMoney } from "@/shared/lib/money";
import { purchaseCatalogQueryOptions } from "../api";
import type { PurchaseCatalog } from "../api/purchase-catalog";

import LucideAlertTriangleRaw from "~icons/lucide/alert-triangle";
import LucideCheckCircle2Raw from "~icons/lucide/check-circle-2";
import LucideClock3Raw from "~icons/lucide/clock-3";
import LucideCopyRaw from "~icons/lucide/copy";
import LucideFileTextRaw from "~icons/lucide/file-text";
import LucideQrCodeRaw from "~icons/lucide/qr-code";
import LucideShieldCheckRaw from "~icons/lucide/shield-check";
import LucideShoppingBagRaw from "~icons/lucide/shopping-bag";
import LucideTicketRaw from "~icons/lucide/ticket";
import LucideXCircleRaw from "~icons/lucide/x-circle";

type IconComponent = Component<{
  class?: string;
}>;

const asIcon = (icon: unknown) => icon as IconComponent;

const LucideAlertTriangle = asIcon(LucideAlertTriangleRaw);
const LucideCheckCircle2 = asIcon(LucideCheckCircle2Raw);
const LucideClock3 = asIcon(LucideClock3Raw);
const LucideCopy = asIcon(LucideCopyRaw);
const LucideFileText = asIcon(LucideFileTextRaw);
const LucideQrCode = asIcon(LucideQrCodeRaw);
const LucideShieldCheck = asIcon(LucideShieldCheckRaw);
const LucideShoppingBag = asIcon(LucideShoppingBagRaw);
const LucideTicket = asIcon(LucideTicketRaw);
const LucideXCircle = asIcon(LucideXCircleRaw);

const statusCopy = {
  pending: [
    "Escaneie para pagar",
    "Confirme o pagamento pelo app do seu banco. A liberação é automática.",
    LucideClock3,
  ],
  approved: [
    "Pedido confirmado",
    "Seu pagamento foi aprovado e os itens já estão disponíveis na sua conta.",
    LucideCheckCircle2,
  ],
  expired: [
    "Reserva expirada",
    "O prazo para pagamento acabou e os itens foram liberados para outros compradores.",
    LucideAlertTriangle,
  ],
  cancelled: [
    "Pedido cancelado",
    "Nenhuma cobrança adicional será feita nesse pedido.",
    LucideXCircle,
  ],
  failed: [
    "Algo deu errado",
    "Não conseguimos processar o pagamento agora.",
    LucideXCircle,
  ],
  rejected: [
    "Pagamento recusado",
    "Não conseguimos processar o seu pagamento no momento.",
    LucideXCircle,
  ],
  refunded: [
    "Pedido reembolsado",
    "O pagamento foi devolvido e os itens desta compra foram cancelados.",
    LucideCheckCircle2,
  ],
} as const;

export default function CheckoutPage(props: { purchase: Purchase }) {
  const queryClient = useQueryClient();

  const catalogData = createMemo(() =>
    queryClient
      .fetchQuery(purchaseCatalogQueryOptions(props.purchase))
      .catch((error) => {
        console.error("Erro ao carregar catálogo da compra:", error);

        return {} as PurchaseCatalog;
      }),
  );

  const [copied, setCopied] = createSignal(false);
  const [qrImage, setQrImage] = createSignal<string>();

  const statusData = createMemo(
    () => statusCopy[props.purchase.status],
  );

  const statusHeading = () => statusData()[0];
  const statusBody = () => statusData()[1];

  const renderStatusIcon = () => {
    const StatusIcon = statusData()[2];

    return <StatusIcon class="relative z-10 size-12" />;
  };

  const subtotal = createMemo(() =>
    props.purchase.items.reduce(
      (sum, item) =>
        sum + item.unit_price_cents * item.quantity,
      0,
    ),
  );

  const isPending = createMemo(
    () => props.purchase.status === "pending",
  );

  const isPixPending = createMemo(
    () =>
      isPending() &&
      props.purchase.payment_method === "pix",
  );

  const heading = createMemo(() =>
    isPixPending()
      ? statusHeading()
      : isPending()
        ? "Pagamento em processamento"
        : statusHeading(),
  );

  const body = createMemo(() =>
    isPixPending()
      ? statusBody()
      : isPending()
        ? "Aguarde a confirmação do pagamento. Você será atualizado automaticamente."
        : statusBody(),
  );

  const tone = createMemo(() =>
    props.purchase.status === "approved"
      ? "text-emerald-600 bg-emerald-500/15"
      : isPending()
        ? "text-amber-600 bg-amber-500/15"
        : "text-muted-foreground bg-muted",
  );

  createEffect(
    () => {
      if (
        !isPixPending() ||
        props.purchase.qr_code_base64
      ) {
        return undefined;
      }

      return props.purchase.qr_code ?? undefined;
    },

    (qrCode) => {
      if (!qrCode) {
        setQrImage(undefined);
        return;
      }

      let cancelled = false;

      QRCode.toDataURL(qrCode, {
        margin: 0,
        width: 220,
      })
        .then((image) => {
          if (!cancelled) setQrImage(image);
        })
        .catch(() => {
          if (!cancelled) setQrImage(undefined);
        });

      return () => {
        cancelled = true;
      };
    },
  );

  return (
    <main class="min-h-screen min-w-0 overflow-x-hidden bg-background px-2 py-6 text-foreground sm:px-8 sm:py-14 min-[360px]:px-4 pb-28!">
      <div class="relative mx-auto grid min-w-0 w-full max-w-4xl overflow-visible rounded-xl border border-dashed border-border bg-card shadow-xl md:grid-cols-[360px_1fr]">
        {/* Order Summary                */}
        <section class="flex min-w-0 flex-col overflow-hidden rounded-t-[calc(0.75rem-2px)] bg-muted p-4 sm:p-8 min-[360px]:p-6 md:rounded-l-[calc(0.75rem-2px)] md:rounded-tr-none md:rounded-br-none">
          <p class="text-[11px] font-bold uppercase tracking-widest text-muted-foreground">
            Resumo do pedido
          </p>

          <div class="mt-6 space-y-4">
            <Loading
              fallback={
                <div class="flex min-h-20 items-center justify-center text-sm text-muted-foreground">
                  Carregando itens...
                </div>
              }
            >
              <For each={props.purchase.items}>
                {(item) => {
                  const itemKey =
                    `${item.item_type}:${item.item_id}`;

                  const itemData = () =>
                    catalogData()[itemKey];

                  return (
                    <div class="flex min-w-0 items-center gap-2 min-[360px]:gap-3">
                      <div class="relative shrink-0">
                        <Show
                          when={itemData()?.image}
                          fallback={
                            <div class="flex size-14 items-center justify-center rounded-xl border border-border bg-background text-muted-foreground">
                              <Show
                                when={
                                  item.item_type === "ticket"
                                }
                                fallback={
                                  <LucideQrCode class="size-5" />
                                }
                              >
                                <LucideTicket class="size-5" />
                              </Show>
                            </div>
                          }
                        >
                          <img
                            src={
                              itemData()?.image ??
                              undefined
                            }
                            alt=""
                            width={56}
                            height={56}
                            class="size-14 rounded-xl object-cover"
                          />
                        </Show>

                        <span
                          class="absolute -bottom-2 -right-2 flex size-7 items-center justify-center rounded-full bg-primary text-primary-foreground ring-4 ring-muted"
                          title={item.item_type}
                        >
                          <Show
                            when={
                              item.item_type === "ticket"
                            }
                            fallback={
                              <Show
                                when={
                                  item.item_type ===
                                  "product"
                                }
                                fallback={
                                  <LucideFileText class="size-3.5" />
                                }
                              >
                                <LucideShoppingBag class="size-3.5" />
                              </Show>
                            }
                          >
                            <LucideTicket class="size-3.5" />
                          </Show>
                        </span>
                      </div>

                      <div class="min-w-0 flex-1">
                        <p class="truncate text-sm font-semibold">
                          {itemData()?.name ??
                            "Item não identificado"}
                        </p>

                        <p class="text-xs text-muted-foreground">
                          {itemData()?.description ??
                            `${item.quantity}x unidade`}
                        </p>
                      </div>

                      <span class="shrink-0 text-xs font-semibold tabular-nums min-[360px]:text-sm">
                        {formatMoney(
                          item.unit_price_cents *
                          item.quantity,
                          props.purchase.currency,
                        )}
                      </span>
                    </div>
                  );
                }}
              </For>
            </Loading>
          </div>

          <div class="mt-6 space-y-2 border-t border-border pt-5 text-sm">
            <div class="flex justify-between">
              <span class="text-muted-foreground">
                Subtotal
              </span>

              <span>
                {formatMoney(
                  subtotal(),
                  props.purchase.currency,
                )}
              </span>
            </div>
          </div>

          <div class="mt-4 flex justify-between border-t border-border pt-4 font-bold">
            <span>Total</span>

            <span class="text-xl tabular-nums">
              {formatMoney(
                props.purchase.total_cents,
                props.purchase.currency,
              )}
            </span>
          </div>

          <div class="mt-auto flex items-center gap-1.5 pt-8 text-xs text-muted-foreground">
            <LucideShieldCheck class="size-3.5" />
            Checkout protegido
          </div>
        </section>

        {/* Payment Status             */}
        <section class="relative flex min-w-0 flex-col items-center justify-center overflow-visible rounded-b-[calc(0.75rem-2px)] border-t border-dashed border-border px-4 py-8 text-center sm:px-10 sm:py-10 md:rounded-r-[calc(0.75rem-2px)] md:rounded-bl-none md:rounded-tl-none md:border-t-0 md:border-l min-[360px]:px-8">
          <span class="pointer-events-none z-50 absolute -top-3.5 left-1/2 size-7 -translate-x-1/2 rounded-full bg-background md:-left-3.5 md:top-1/2 md:-translate-y-1/2 md:translate-x-0" />

          <div
            class={`relative flex size-28 items-center justify-center rounded-full ${tone()}`}
          >
            <span class="absolute -inset-4 z-0 rounded-full bg-current opacity-10 blur-xl" />

            {renderStatusIcon()}
          </div>

          <h1 class="mt-5 min-h-16 max-w-65 text-xl font-extrabold leading-tight tracking-tight min-[360px]:text-[22px]">
            {heading()}
          </h1>

          <div class="mt-1.5 min-h-10 max-w-xs text-sm leading-relaxed text-muted-foreground">
            <p>
              {body()}

              {!isPending() &&
                props.purchase.status_reason
                ? ` ${props.purchase.status_reason}.`
                : ""}
            </p>
          </div>

          {/* PIX                             */}
          <div class="mt-4 flex w-full max-w-xs items-center justify-center">
            <Show
              when={isPixPending()}
              fallback={
                <Show
                  when={isPending()}
                  fallback={
                    <div class="min-h-10" />
                  }
                >
                  <p class="text-sm text-muted-foreground">
                    Pagamento em processamento.
                  </p>
                </Show>
              }
            >
              <div class="w-full space-y-3">
                <div class="mx-auto flex size-36 items-center justify-center rounded-2xl border border-border p-3 min-[360px]:size-44">
                  <Show
                    when={
                      props.purchase.qr_code_base64
                    }
                    fallback={
                      <Show
                        when={qrImage()}
                        fallback={
                          <LucideQrCode class="size-full text-foreground" />
                        }
                      >
                        <img
                          src={qrImage()}
                          alt="QR Code Pix"
                          width={176}
                          height={176}
                          class="size-full"
                        />
                      </Show>
                    }
                  >
                    <img
                      src={`data:image/png;base64,${props.purchase.qr_code_base64}`}
                      alt="QR Code Pix"
                      width={176}
                      height={176}
                      class="size-full"
                    />
                  </Show>
                </div>

                <div class="flex items-center justify-center gap-1.5 text-xs text-muted-foreground">
                  <LucideClock3 class="size-3.5" />
                  Aguardando confirmação do pagamento
                </div>

                <Show when={props.purchase.qr_code}>
                  <button
                    type="button"
                    class="flex min-w-0 w-full items-center gap-2 rounded-lg border border-border px-3 py-2 text-left text-xs text-muted-foreground hover:bg-muted"
                    onClick={() => {
                      navigator.clipboard?.writeText(
                        props.purchase.qr_code ?? "",
                      );

                      setCopied(true);

                      setTimeout(
                        () => setCopied(false),
                        1800,
                      );
                    }}
                  >
                    <span class="min-w-0 flex-1 truncate">
                      {props.purchase.qr_code}
                    </span>

                    <Show
                      when={copied()}
                      fallback={
                        <>
                          <LucideCopy class="size-3.5 shrink-0" />
                          Copiar
                        </>
                      }
                    >
                      Copiado
                    </Show>
                  </button>
                </Show>
              </div>
            </Show>
          </div>

          {/* Purchase Data                */}
          <div class="mt-2 w-full max-w-xs rounded-xl border border-border bg-muted/50 p-4 text-left text-xs">
            <div class="flex min-w-0 items-center justify-between gap-2">
              <span class="min-w-0 text-muted-foreground">
                Forma de pagamento
              </span>

              <span class="shrink-0 text-right font-semibold">
                {props.purchase.payment_method === "pix"
                  ? "Pix"
                  : (props.purchase.payment_method ??
                    "Não informado")}
              </span>
            </div>

            <div class="mt-2 flex min-w-0 items-center justify-between gap-2">
              <span class="min-w-0 text-muted-foreground">
                Pedido realizado em
              </span>

              <span class="shrink-0 text-right font-semibold">
                {new Date(
                  props.purchase.created_at ??
                  props.purchase.expires_at,
                ).toLocaleDateString("pt-BR")}
              </span>
            </div>

            <div class="mt-2 flex min-w-0 items-center justify-between gap-2">
              <span class="min-w-0 text-muted-foreground">
                Pedido
              </span>

              <span class="shrink-0 font-mono font-semibold">
                #
                {props.purchase.purchase_id.slice(
                  0,
                  8,
                )}
              </span>
            </div>
          </div>

          {/* Go Back                         */}
          <Show when={!isPending()}>
            <Link
              to="/profile"
              search={{
                tab: "purchases",
              }}
              class="mt-6 inline-flex min-h-12 w-full max-w-xs items-center justify-center rounded-xl bg-primary px-4 py-3 text-sm font-semibold text-primary-foreground hover:bg-primary/90"
            >
              Voltar para minhas compras
            </Link>
          </Show>
        </section>
      </div>
    </main>
  );
}