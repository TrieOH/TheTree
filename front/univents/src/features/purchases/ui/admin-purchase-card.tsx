import type { EditionPurchase } from "@trieoh/univents-api/schemas";
import type { LucideIcon } from "lucide-react";
import {
  CalendarClock,
  ChevronDown,
  CreditCard,
  Hash,
  Mail,
  Package,
  Ticket,
  Users,
} from "lucide-react";
import { useState } from "react";

const STATUS_STYLES: Record<string, StatusStyle> = {
  paid: {
    label: "Pago",
    badge: "bg-emerald-950 text-emerald-400 border border-emerald-900",
    dot: "bg-emerald-400",
  },
  completed: {
    label: "Concluído",
    badge: "bg-emerald-950 text-emerald-400 border border-emerald-900",
    dot: "bg-emerald-400",
  },
  confirmed: {
    label: "Confirmado",
    badge: "bg-emerald-950 text-emerald-400 border border-emerald-900",
    dot: "bg-emerald-400",
  },
  approved: {
    label: "Aprovado",
    badge: "bg-emerald-950 text-emerald-400 border border-emerald-900",
    dot: "bg-emerald-400",
  },
  expired: {
    label: "Expirado",
    badge: "bg-slate-800 text-slate-400 border border-slate-700",
    dot: "bg-slate-400",
  },
  rejected: {
    label: "Rejeitado",
    badge: "bg-rose-950 text-rose-400 border border-rose-900",
    dot: "bg-rose-400",
  },
  pending: {
    label: "Pendente",
    badge: "bg-amber-950 text-amber-400 border border-amber-900",
    dot: "bg-amber-400",
  },
  processing: {
    label: "Processando",
    badge: "bg-amber-950 text-amber-400 border border-amber-900",
    dot: "bg-amber-400",
  },
  failed: {
    label: "Falhou",
    badge: "bg-rose-950 text-rose-400 border border-rose-900",
    dot: "bg-rose-400",
  },
  cancelled: {
    label: "Cancelado",
    badge: "bg-rose-950 text-rose-400 border border-rose-900",
    dot: "bg-rose-400",
  },
  refunded: {
    label: "Reembolsado",
    badge: "bg-slate-800 text-slate-400 border border-slate-700",
    dot: "bg-slate-400",
  },
};

const DEFAULT_STATUS: StatusStyle = {
  badge: "bg-slate-800 text-slate-400 border border-slate-700",
  dot: "bg-slate-400",
};

const REFUNDABLE_STATUSES = ["paid", "completed", "confirmed", "approved"];

const STATUS_MESSAGES: Record<string, string> = {
  approved: "Pagamento aprovado.",
  cancelled: "Compra cancelada.",
  completed: "Compra concluída.",
  confirmed: "Compra confirmada.",
  expired: "O prazo desta compra expirou.",
  failed: "Não foi possível concluir o pagamento.",
  paid: "Pagamento confirmado.",
  pending: "Aguardando confirmação do pagamento.",
  processing: "Pagamento em processamento.",
  refunded: "Compra reembolsada.",
  rejected: "Pagamento rejeitado.",
};

const ITEM_META: Record<string, ItemMeta> = {
  ticket: { icon: Ticket, label: "Ingresso" },
  product: { icon: Package, label: "Produto" },
  program_occurrence: { icon: CalendarClock, label: "Sessão" },
};

interface StatusStyle {
  label?: string;
  badge: string;
  dot: string;
}

interface ItemMeta {
  icon: LucideIcon;
  label: string;
}

type PurchaseItem = EditionPurchase["items"][number];

interface AdminPurchaseCardProps {
  purchase: EditionPurchase;
  onRefund?: (purchase: EditionPurchase) => void;
}

function statusStyle(status: EditionPurchase["status"]): Required<StatusStyle> {
  const key = status.toLowerCase();
  const style = STATUS_STYLES[key] ?? DEFAULT_STATUS;
  return {
    badge: style.badge,
    dot: style.dot,
    label: style.label ?? humanizeStatus(status),
  };
}

function humanizeStatus(status: string) {
  return status
    .replaceAll("_", " ")
    .replace(/\b\w/g, (character) => character.toUpperCase());
}

function formatMoney(cents: number, currency: string | null | undefined) {
  try {
    return new Intl.NumberFormat("pt-BR", {
      style: "currency",
      currency: currency || "BRL",
    }).format((cents || 0) / 100);
  } catch {
    return `${((cents || 0) / 100).toFixed(2)} ${currency || ""}`;
  }
}

function formatDate(iso: string | null | undefined) {
  if (!iso) return "—";
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return "—";
  return d.toLocaleDateString("pt-BR", {
    day: "2-digit",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

function shortId(id: string | null | undefined) {
  if (!id) return "—";
  return id.length > 10 ? `${id.slice(0, 8)}…` : id;
}

function shortItemId(id: string | null | undefined) {
  if (!id) return "";
  return id.length > 8 ? id.slice(0, 8) : id;
}

export function AdminPurchaseCard({
  purchase,
  onRefund,
}: AdminPurchaseCardProps) {
  const [open, setOpen] = useState(false);
  const s = statusStyle(purchase.status);
  const attendees = purchase.attendees || [];
  const items = purchase.items || [];
  const isRefundable = REFUNDABLE_STATUSES.includes(
    String(purchase.status || "").toLowerCase(),
  );

  return (
    <div className="flex w-full min-w-0 flex-col overflow-hidden rounded-xl border border-border bg-muted text-card-foreground shadow-lg xl:flex-row">
      {/* Coluna esquerda: identidade + valor */}
      <div className="flex min-w-0 flex-col gap-2.5 p-4">
        <div className="flex items-start justify-between gap-3">
          <div>
            <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
              Compra
            </p>
            <p className="font-mono flex items-center gap-1.5 text-xs mt-0.5 text-foreground">
              <Hash size={12} className="text-muted-foreground" />
              {shortId(purchase.purchase_id)}
            </p>
          </div>
          <span
            className={`inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-medium ${s.badge}`}
          >
            <span className={`h-1.5 w-1.5 rounded-full ${s.dot}`} />
            {s.label}
          </span>
        </div>

        <p className="min-h-8 text-xs leading-relaxed text-muted-foreground">
          {purchase.status_reason ??
            STATUS_MESSAGES[purchase.status.toLowerCase()] ??
            `Status: ${humanizeStatus(purchase.status)}.`}
        </p>

        <div className="space-y-1 text-xs text-muted-foreground">
          <div className="flex items-center gap-1.5">
            <CalendarClock
              size={13}
              className="text-muted-foreground shrink-0"
            />
            <span>{formatDate(purchase.created_at)}</span>
          </div>
          {purchase.payer_email && (
            <div className="flex items-center gap-1.5">
              <Mail size={13} className="text-muted-foreground shrink-0" />
              <span className="truncate">{purchase.payer_email}</span>
            </div>
          )}
          {purchase.payment_method && (
            <div className="flex items-center gap-1.5">
              <CreditCard
                size={13}
                className="text-muted-foreground shrink-0"
              />
              <span className="capitalize">
                {purchase.payment_method.replace(/_/g, " ")}
              </span>
            </div>
          )}
        </div>

        <div className="mt-auto pt-1">
          <p className="text-2xl font-semibold leading-none text-amber-400">
            {formatMoney(purchase.total_cents, purchase.currency)}
          </p>
          <p className="text-xs mt-1 text-muted-foreground">
            {items.length} {items.length === 1 ? "item" : "itens"}
          </p>

          <button
            type="button"
            disabled={!isRefundable}
            onClick={() => isRefundable && onRefund?.(purchase)}
            className="mt-3 w-full rounded-lg border border-destructive/30 bg-background px-3 py-1.5 text-xs font-medium text-destructive transition-colors hover:bg-destructive/10 disabled:cursor-not-allowed disabled:border-border disabled:text-muted-foreground disabled:hover:bg-background"
          >
            {isRefundable ? "Reembolsar" : "Indisponível"}
          </button>
        </div>
      </div>

      <div className="relative shrink-0 w-full">
        <div className="absolute -top-1 -left-1.5 h-3 w-3 rounded-full bg-background" />
        <div className="absolute -top-1 -right-1.5 h-3 w-3 rounded-full bg-background" />
        <div className="w-full h-px border-t border-dashed border-border" />
      </div>

      <div className="flex min-w-0 flex-1 flex-col gap-2.5 p-4">
        <div>
          <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground mb-1.5">
            Itens
          </p>
          <ul className="space-y-1">
            {items.map((item: PurchaseItem, i: number) => {
              const meta = ITEM_META[item.item_type] || {
                icon: Package,
                label: item.item_type,
              };
              const Icon = meta.icon;
              const subtotal =
                (item.unit_price_cents || 0) * (item.quantity || 1);
              return (
                <li
                  key={`${item.item_id}-${i}`}
                  className="flex items-center justify-between gap-2 text-xs"
                >
                  <span className="flex items-center gap-1.5 min-w-0 text-foreground">
                    <Icon size={13} className="text-amber-500 shrink-0" />
                    <span className="truncate">{meta.label}</span>
                    <span className="font-mono truncate text-muted-foreground">
                      #{shortItemId(item.item_id)}
                    </span>
                  </span>
                  <span className="font-mono shrink-0 text-muted-foreground">
                    {item.quantity}×{" "}
                    {formatMoney(item.unit_price_cents, purchase.currency)}
                    {item.quantity > 1 && (
                      <span className="text-muted-foreground">
                        {" "}
                        = {formatMoney(subtotal, purchase.currency)}
                      </span>
                    )}
                  </span>
                </li>
              );
            })}
            {items.length === 0 && (
              <li className="text-xs italic text-muted-foreground">
                Nenhum item
              </li>
            )}
          </ul>
        </div>

        <button
          type="button"
          onClick={() => setOpen((v) => !v)}
          className="mt-auto flex items-center justify-between rounded-lg border border-border bg-muted px-2.5 py-1.5 text-left transition-colors hover:bg-accent"
        >
          <span className="flex items-center gap-2 text-xs text-foreground">
            <Users size={13} className="text-muted-foreground" />
            {attendees.length}{" "}
            {attendees.length === 1 ? "participante" : "participantes"}
          </span>
          <ChevronDown
            size={15}
            className={`text-muted-foreground transition-transform duration-200 ${open ? "rotate-180" : ""}`}
          />
        </button>

        {open && (
          <ul className="space-y-1 pl-1">
            {attendees.map((attendee, i) => (
              <li
                key={`${attendee.email}-${i}`}
                className="flex items-center justify-between gap-3 text-xs"
              >
                <span className="text-foreground">{attendee.name}</span>
                <span className="font-mono truncate text-muted-foreground">
                  {attendee.email}
                </span>
              </li>
            ))}
            {attendees.length === 0 && (
              <li className="text-xs italic text-muted-foreground">
                Nenhum participante
              </li>
            )}
          </ul>
        )}
      </div>
    </div>
  );
}
