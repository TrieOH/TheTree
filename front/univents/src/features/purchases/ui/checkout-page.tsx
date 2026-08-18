import { useQuery } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import type { Purchase } from "@trieoh/univents-api/schemas";
import {
  AlertTriangle,
  CheckCircle2,
  Clock3,
  Copy,
  FileText,
  QrCode,
  ShieldCheck,
  ShoppingBag,
  Ticket,
  XCircle,
} from "lucide-react";
import QRCode from "qrcode";
import { useEffect, useState } from "react";
import { resolvePurchaseCatalog } from "../api/purchase-catalog";

const money = (cents: number, currency: string) =>
  new Intl.NumberFormat("pt-BR", { style: "currency", currency }).format(
    cents / 100,
  );

const statusCopy = {
  pending: [
    "Escaneie para pagar",
    "Confirme o pagamento pelo app do seu banco. A liberação é automática.",
    Clock3,
  ],
  approved: [
    "Pedido confirmado",
    "Seu pagamento foi aprovado e os itens já estão disponíveis na sua conta.",
    CheckCircle2,
  ],
  expired: [
    "Reserva expirada",
    "O prazo para pagamento acabou e os itens foram liberados para outros compradores.",
    AlertTriangle,
  ],
  cancelled: [
    "Pedido cancelado",
    "Nenhuma cobrança adicional será feita nesse pedido.",
    XCircle,
  ],
  failed: [
    "Algo deu errado",
    "Não conseguimos processar o pagamento agora.",
    XCircle,
  ],
  rejected: [
    "Pagamento recusado",
    "Não conseguimos processar o seu pagamento no momento.",
    XCircle,
  ],
} as const;

export default function CheckoutPage({ purchase }: { purchase: Purchase }) {
  const { data: catalog = {} } = useQuery({
    queryKey: ["purchase-catalog", purchase.purchase_id],
    queryFn: () =>
      resolvePurchaseCatalog(
        purchase.items.map((item) => ({
          ...item,
          edition_id: purchase.edition_id,
        })),
      ),
  });
  const [copied, setCopied] = useState(false);
  const [qrImage, setQrImage] = useState<string>();
  const [statusHeading, statusBody, StatusIcon] = statusCopy[purchase.status];
  const subtotal = purchase.items.reduce(
    (sum, item) => sum + item.unit_price_cents * item.quantity,
    0,
  );
  const isPending = purchase.status === "pending";
  const isPixPending = isPending && purchase.payment_method === "pix";
  const heading = isPixPending
    ? statusHeading
    : isPending
      ? "Pagamento em processamento"
      : statusHeading;
  const body = isPixPending
    ? statusBody
    : isPending
      ? "Aguarde a confirmação do pagamento. Você será atualizado automaticamente."
      : statusBody;
  const tone =
    purchase.status === "approved"
      ? "text-emerald-600 bg-emerald-500/15"
      : isPending
        ? "text-amber-600 bg-amber-500/15"
        : "text-muted-foreground bg-muted";

  useEffect(() => {
    if (!isPixPending || !purchase.qr_code || purchase.qr_code_base64) {
      setQrImage(undefined);
      return;
    }
    void QRCode.toDataURL(purchase.qr_code, { margin: 0, width: 220 })
      .then(setQrImage)
      .catch(() => setQrImage(undefined));
  }, [isPixPending, purchase.qr_code, purchase.qr_code_base64]);

  return (
    <main className="min-h-screen min-w-0 overflow-x-hidden bg-background px-2 py-6 text-foreground sm:px-8 sm:py-14 min-[360px]:px-4 pb-28!">
      <div className="relative mx-auto grid min-w-0 w-full max-w-4xl overflow-visible rounded-xl border border-dashed border-border bg-card shadow-xl md:grid-cols-[360px_1fr]">
        <section className="flex min-w-0 flex-col overflow-hidden rounded-t-[calc(0.75rem-2px)] bg-muted p-4 sm:p-8 min-[360px]:p-6 md:rounded-l-[calc(0.75rem-2px)] md:rounded-tr-none md:rounded-br-none">
          <p className="text-[11px] font-bold uppercase tracking-widest text-muted-foreground">
            Resumo do pedido
          </p>
          <div className="mt-6 space-y-4">
            {purchase.items.map((item) => (
              <div
                key={`${item.item_type}:${item.item_id}`}
                className="flex min-w-0 items-center gap-2 min-[360px]:gap-3"
              >
                <div className="relative shrink-0">
                  {catalog[`${item.item_type}:${item.item_id}`]?.image ? (
                    <img
                      src={
                        catalog[`${item.item_type}:${item.item_id}`]?.image ??
                        undefined
                      }
                      alt=""
                      className="size-14 rounded-xl object-cover"
                    />
                  ) : (
                    <div className="flex size-14 items-center justify-center rounded-xl border border-border bg-background text-muted-foreground">
                      {item.item_type === "ticket" ? (
                        <Ticket className="size-5" />
                      ) : (
                        <QrCode className="size-5" />
                      )}
                    </div>
                  )}
                  <span
                    className="absolute -bottom-2 -right-2 flex size-7 items-center justify-center rounded-full bg-primary text-primary-foreground ring-4 ring-muted"
                    title={item.item_type}
                  >
                    {item.item_type === "ticket" ? (
                      <Ticket className="size-3.5" />
                    ) : item.item_type === "product" ? (
                      <ShoppingBag className="size-3.5" />
                    ) : (
                      <FileText className="size-3.5" />
                    )}
                  </span>
                </div>
                <div className="min-w-0 flex-1">
                  <p className="truncate text-sm font-semibold">
                    {catalog[`${item.item_type}:${item.item_id}`]?.name ??
                      "Item não identificado"}
                  </p>
                  <p className="text-xs text-muted-foreground">
                    {catalog[`${item.item_type}:${item.item_id}`]
                      ?.description ?? `${item.quantity}x unidade`}
                  </p>
                </div>
                <span className="shrink-0 text-xs font-semibold tabular-nums min-[360px]:text-sm">
                  {money(
                    item.unit_price_cents * item.quantity,
                    purchase.currency,
                  )}
                </span>
              </div>
            ))}
          </div>

          <div className="mt-6 space-y-2 border-t border-border pt-5 text-sm">
            <div className="flex justify-between">
              <span className="text-muted-foreground">Subtotal</span>
              <span>{money(subtotal, purchase.currency)}</span>
            </div>
          </div>
          <div className="mt-4 flex justify-between border-t border-border pt-4 font-bold">
            <span>Total</span>
            <span className="text-xl tabular-nums">
              {money(purchase.total_cents, purchase.currency)}
            </span>
          </div>
          <div className="mt-auto flex items-center gap-1.5 pt-8 text-xs text-muted-foreground">
            <ShieldCheck className="size-3.5" /> Checkout protegido
          </div>
        </section>

        <section className="relative flex min-w-0 flex-col items-center justify-center overflow-visible rounded-b-[calc(0.75rem-2px)] border-t border-dashed border-border px-4 py-8 text-center sm:px-10 sm:py-10 md:rounded-r-[calc(0.75rem-2px)] md:rounded-bl-none md:rounded-tl-none md:border-t-0 md:border-l min-[360px]:px-8">
          <span className="pointer-events-none z-50 absolute -top-3.5 left-1/2 size-7 -translate-x-1/2 rounded-full bg-background md:-left-3.5 md:top-1/2 md:-translate-y-1/2 md:translate-x-0" />
          <div
            className={`relative flex size-28 items-center justify-center rounded-full ${tone}`}
          >
            <span className="absolute -inset-4 z-0 rounded-full bg-current opacity-10 blur-xl" />
            <StatusIcon className="relative z-10 size-12" />
          </div>
          <h1 className="mt-5 min-h-16 max-w-65 text-xl font-extrabold leading-tight tracking-tight min-[360px]:text-[22px]">
            {heading}
          </h1>
          <div className="mt-1.5 min-h-10 max-w-xs text-sm leading-relaxed text-muted-foreground">
            <p>
              {body}
              {!isPending && purchase.status_reason
                ? ` ${purchase.status_reason}.`
                : ""}
            </p>
          </div>
          <div className="mt-4 flex w-full max-w-xs items-center justify-center">
            {isPixPending ? (
              <div className="w-full space-y-3">
                <div className="mx-auto flex size-36 items-center justify-center rounded-2xl border border-border p-3 min-[360px]:size-44">
                  {purchase.qr_code_base64 ? (
                    <img
                      src={`data:image/png;base64,${purchase.qr_code_base64}`}
                      alt="QR Code Pix"
                      className="size-full"
                    />
                  ) : qrImage ? (
                    <img
                      src={qrImage}
                      alt="QR Code Pix"
                      className="size-full"
                    />
                  ) : (
                    <QrCode className="size-full text-foreground" />
                  )}
                </div>
                <div className="flex items-center justify-center gap-1.5 text-xs text-muted-foreground">
                  <Clock3 className="size-3.5" /> Aguardando confirmação do
                  pagamento
                </div>
                {purchase.qr_code && (
                  <button
                    type="button"
                    className="flex min-w-0 w-full items-center gap-2 rounded-lg border border-border px-3 py-2 text-left text-xs text-muted-foreground hover:bg-muted"
                    onClick={() => {
                      navigator.clipboard?.writeText(purchase.qr_code ?? "");
                      setCopied(true);
                      setTimeout(() => setCopied(false), 1800);
                    }}
                  >
                    <span className="min-w-0 flex-1 truncate">
                      {purchase.qr_code}
                    </span>
                    {copied ? (
                      "Copiado"
                    ) : (
                      <>
                        <Copy className="size-3.5 shrink-0" /> Copiar
                      </>
                    )}
                  </button>
                )}
              </div>
            ) : isPending ? (
              <p className="text-sm text-muted-foreground">
                Pagamento em processamento.
              </p>
            ) : (
              <div className="min-h-10" />
            )}
          </div>
          <div className="mt-2 w-full max-w-xs rounded-xl border border-border bg-muted/50 p-4 text-left text-xs">
            <div className="flex min-w-0 items-center justify-between gap-2">
              <span className="min-w-0 text-muted-foreground">
                Forma de pagamento
              </span>
              <span className="shrink-0 text-right font-semibold">
                {purchase.payment_method === "pix"
                  ? "Pix"
                  : (purchase.payment_method ?? "Não informado")}
              </span>
            </div>
            <div className="mt-2 flex min-w-0 items-center justify-between gap-2">
              <span className="min-w-0 text-muted-foreground">
                Pedido realizado em
              </span>
              <span className="shrink-0 text-right font-semibold">
                {new Date(
                  purchase.created_at ?? purchase.expires_at,
                ).toLocaleDateString("pt-BR")}
              </span>
            </div>
            <div className="mt-2 flex min-w-0 items-center justify-between gap-2">
              <span className="min-w-0 text-muted-foreground">Pedido</span>
              <span className="shrink-0 font-mono font-semibold">
                #{purchase.purchase_id.slice(0, 8)}
              </span>
            </div>
          </div>
          {!isPending && (
            <Link
              to="/profile"
              search={{ tab: "purchases" }}
              className="mt-6 inline-flex min-h-12 w-full max-w-xs items-center justify-center rounded-xl bg-primary px-4 py-3 text-sm font-semibold text-primary-foreground hover:bg-primary/90"
            >
              Voltar para minhas compras
            </Link>
          )}
        </section>
      </div>
    </main>
  );
}
