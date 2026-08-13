import type { Purchase } from "@trieoh/univents-api/schemas";
import { cn } from "@/shared/lib/utils";

const statusLabel = {
  pending: "Aguardando pagamento",
  approved: "Aprovada",
  expired: "Expirada",
  cancelled: "Cancelada",
} satisfies Record<Purchase["status"], string>;

const itemTypeLabel = {
  ticket: "Ingresso",
  product: "Produto",
  program_occurrence: "Atividade",
} satisfies Record<Purchase["items"][number]["item_type"], string>;

const formatCurrency = (cents: number, currency: string) =>
  new Intl.NumberFormat("pt-BR", { style: "currency", currency }).format(
    cents / 100,
  );

export function PurchaseSummary({ purchase }: { purchase: Purchase }) {
  const hasPixQr =
    purchase.status === "pending" &&
    purchase.payment_method === "pix" &&
    purchase.qr_code &&
    purchase.qr_code_base64;

  return (
    <article className="space-y-4 border border-border bg-card p-5">
      <header className="flex flex-wrap items-start justify-between gap-3">
        <div>
          <p className="text-xs text-muted-foreground">
            Compra {purchase.purchase_id.slice(0, 8)}
          </p>
          <p className="text-sm font-medium">
            {new Intl.DateTimeFormat("pt-BR", {
              dateStyle: "medium",
              timeStyle: "short",
            }).format(new Date(purchase.created_at ?? purchase.expires_at))}
          </p>
        </div>
        <span
          className={cn(
            "rounded-full px-3 py-1 text-xs font-semibold",
            purchase.status === "approved" &&
              "bg-emerald-500/10 text-emerald-700 dark:text-emerald-400",
            purchase.status === "pending" &&
              "bg-amber-500/10 text-amber-700 dark:text-amber-400",
            (purchase.status === "expired" ||
              purchase.status === "cancelled") &&
              "bg-muted text-muted-foreground",
          )}
        >
          {statusLabel[purchase.status]}
        </span>
      </header>

      <div className="divide-y divide-border">
        {purchase.items.map((item) => (
          <div
            key={`${item.item_type}:${item.item_id}`}
            className="flex items-center justify-between gap-4 py-3 text-sm"
          >
            <div>
              <p className="font-medium">{itemTypeLabel[item.item_type]}</p>
              <p className="text-xs text-muted-foreground">
                {item.quantity} ×{" "}
                {formatCurrency(item.unit_price_cents, purchase.currency)}
              </p>
            </div>
            <span className="font-medium tabular-nums">
              {formatCurrency(
                item.unit_price_cents * item.quantity,
                purchase.currency,
              )}
            </span>
          </div>
        ))}
      </div>

      <footer className="flex items-center justify-between border-t border-border pt-4">
        <span className="font-semibold">Total</span>
        <span className="text-lg font-bold text-primary tabular-nums">
          {formatCurrency(purchase.total_cents, purchase.currency)}
        </span>
      </footer>

      {hasPixQr && (
        <section className="flex flex-col items-center gap-4 border-t border-border pt-4 text-center">
          <div>
            <p className="font-semibold">Pague com Pix</p>
            <p className="text-sm text-muted-foreground">
              Escaneie o QR Code ou copie o código abaixo
            </p>
          </div>

          <div className="rounded-xl border border-border bg-white p-4">
            <img
              src={`data:image/png;base64,${purchase.qr_code_base64}`}
              alt="QR Code Pix"
              className="size-48"
            />
          </div>

          <div className="w-full">
            <p className="mb-1 text-xs text-muted-foreground">
              Código Pix copia e cola
            </p>
            <code className="block truncate rounded-md border border-border bg-muted/40 px-3 py-2 text-left text-xs">
              {purchase.qr_code}
            </code>
          </div>
        </section>
      )}
    </article>
  );
}
