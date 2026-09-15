import type { JSX } from "@solidjs/web";
import { For, Show, createSignal } from "solid-js";
import type { EditionPurchase } from "@trieoh/univents-api/schemas";

import CalendarClockIcon from "~icons/lucide/calendar-clock";
import CheckIcon from "~icons/lucide/check";
import ChevronDownIcon from "~icons/lucide/chevron-down";
import CopyIcon from "~icons/lucide/copy";
import HashIcon from "~icons/lucide/hash";
import MailIcon from "~icons/lucide/mail";
import PackageIcon from "~icons/lucide/package";
import TicketIcon from "~icons/lucide/ticket";
import UsersIcon from "~icons/lucide/users";

import { cn } from "@trieoh/ui-solid";
import { Reveal } from "@/shared/ui/Reveal";

const CalendarClock = CalendarClockIcon as unknown as (props: { class?: string }) => JSX.Element;
const Check = CheckIcon as unknown as (props: { class?: string }) => JSX.Element;
const ChevronDown = ChevronDownIcon as unknown as (props: { class?: string }) => JSX.Element;
const Copy = CopyIcon as unknown as (props: { class?: string }) => JSX.Element;
const Hash = HashIcon as unknown as (props: { class?: string }) => JSX.Element;
const Mail = MailIcon as unknown as (props: { class?: string }) => JSX.Element;
const Package = PackageIcon as unknown as (props: { class?: string }) => JSX.Element;
const Ticket = TicketIcon as unknown as (props: { class?: string }) => JSX.Element;
const Users = UsersIcon as unknown as (props: { class?: string }) => JSX.Element;

interface StatusStyle {
  label?: string;
  badge: string;
  dot: string;
}

interface ItemMeta {
  icon: (props: { class?: string }) => JSX.Element;
  label: string;
}

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

function humanizeStatus(status: string): string {
  return status
    .replaceAll("_", " ")
    .replace(/\b\w/g, (character) => character.toUpperCase());
}

function statusStyle(status: string): Required<StatusStyle> {
  const key = status.toLowerCase();
  const style = STATUS_STYLES[key] ?? DEFAULT_STATUS;
  return {
    badge: style.badge,
    dot: style.dot,
    label: style.label ?? humanizeStatus(status),
  };
}

function formatMoney(cents: number, currency?: string | null): string {
  try {
    return new Intl.NumberFormat("pt-BR", {
      style: "currency",
      currency: currency || "BRL",
    }).format((cents || 0) / 100);
  } catch {
    return `${((cents || 0) / 100).toFixed(2)} ${currency || ""}`;
  }
}

function formatDate(iso: string | null | undefined): string {
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

function shortId(id: string | null | undefined): string {
  if (!id) return "—";
  return id.length > 10 ? `${id.slice(0, 8)}…` : id;
}

function shortItemId(id: string | null | undefined): string {
  if (!id) return "";
  return id.length > 8 ? id.slice(0, 8) : id;
}

export interface AdminPurchaseCardProps {
  purchase: EditionPurchase;
  index?: number;
  animate?: boolean;
  onRefund?: (purchase: EditionPurchase) => void;
}

export function AdminPurchaseCard(props: AdminPurchaseCardProps): JSX.Element {
  const [open, setOpen] = createSignal(false);
  const [copied, setCopied] = createSignal(false);

  const s = () => statusStyle(props.purchase.status);
  const attendees = () => props.purchase.attendees || [];
  const items = () => props.purchase.items || [];
  const isRefundable = () =>
    REFUNDABLE_STATUSES.includes(String(props.purchase.status || "").toLowerCase());

  const copyId = async (e: MouseEvent) => {
    e.stopPropagation();
    try {
      await navigator.clipboard.writeText(props.purchase.purchase_id);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch {
      // Ignored
    }
  };

  return (
    <Reveal delay={(props.index ?? 0) * 0.05} animate={props.animate} class="h-full">
      <div class="flex h-full w-full min-w-0 flex-col justify-between overflow-hidden rounded-xl border border-border bg-card/60 text-card-foreground shadow-xs transition-all duration-200 hover:border-foreground/20 hover:shadow-md">
        <div class="flex min-w-0 flex-col gap-2 p-3">
          <div class="flex items-start justify-between gap-2">
            <div>
              <p class="text-[10px] font-semibold uppercase tracking-widest text-muted-foreground">
                Compra
              </p>
              <button
                type="button"
                onClick={copyId}
                class="group/id font-mono flex items-center gap-1 text-xs mt-0.5 text-foreground transition-colors hover:text-primary cursor-pointer"
                title="Clique para copiar o ID da compra"
              >
                <Hash class="size-3 text-muted-foreground shrink-0" />
                <span>{shortId(props.purchase.purchase_id)}</span>
                <Show when={copied()} fallback={<Copy class="size-2.5 opacity-40 group-hover/id:opacity-100" />}>
                  <Check class="size-2.5 text-emerald-500" />
                </Show>
              </button>
            </div>
            <span
              class={cn(
                "inline-flex items-center gap-1.5 rounded-full px-2 py-0.5 text-[11px] font-medium shrink-0",
                s().badge,
              )}
            >
              <span class={cn("size-1.5 rounded-full", s().dot)} />
              {s().label}
            </span>
          </div>

          {/* Altura fixa de 2 linhas para garantir alinhamento perfeito entre cards */}
          <p class="h-8 line-clamp-2 text-[11px] leading-snug text-muted-foreground">
            {props.purchase.status_reason ??
              STATUS_MESSAGES[props.purchase.status.toLowerCase()] ??
              `Status: ${humanizeStatus(props.purchase.status)}.`}
          </p>

          <div class="space-y-1 text-xs text-muted-foreground">
            <div class="flex items-center justify-between gap-2 h-4">
              <div class="flex items-center gap-1.5 min-w-0">
                <CalendarClock class="size-3 text-muted-foreground shrink-0" />
                <span class="truncate text-[11px]">{formatDate(props.purchase.created_at)}</span>
              </div>
              <Show
                when={props.purchase.payment_method}
                fallback={<span class="text-muted-foreground/40 text-[11px]">—</span>}
              >
                {(method) => (
                  <span class="capitalize shrink-0 font-medium text-foreground/80 text-[11px]">
                    {method() === "pix" ? "Pix" : method().replace(/_/g, " ")}
                  </span>
                )}
              </Show>
            </div>

            <div class="flex items-center gap-1.5 min-w-0 h-4">
              <Mail class="size-3 text-muted-foreground shrink-0" />
              <Show
                when={props.purchase.payer_email}
                fallback={
                  <span class="truncate text-[11px] text-muted-foreground/40 italic">
                    Sem e-mail informado
                  </span>
                }
              >
                {(email) => (
                  <span class="truncate text-[11px]" title={email()}>
                    {email()}
                  </span>
                )}
              </Show>
            </div>
          </div>

          <div class="mt-auto pt-2 flex items-end justify-between gap-2 border-t border-border/40">
            <div>
              <p class="text-xl font-bold leading-none text-amber-400">
                {formatMoney(props.purchase.total_cents, props.purchase.currency)}
              </p>
              <p class="text-[10px] mt-1 text-muted-foreground">
                {items().length} {items().length === 1 ? "item" : "itens"}
              </p>
            </div>

            <Show when={isRefundable()}>
              <button
                type="button"
                onClick={() => props.onRefund?.(props.purchase)}
                class="rounded-md border border-destructive/30 bg-background/80 px-2.5 py-1 text-xs font-medium text-destructive transition-colors hover:bg-destructive/10 cursor-pointer"
              >
                Reembolsar
              </button>
            </Show>
          </div>
        </div>

        {/* Notched perforated ticket divider */}
        <div class="relative shrink-0 w-full">
          <div class="absolute -top-1.5 -left-1.5 size-3 rounded-full bg-background" />
          <div class="absolute -top-1.5 -right-1.5 size-3 rounded-full bg-background" />
          <div class="w-full h-px border-t border-dashed border-border" />
        </div>

        <div class="flex min-w-0 flex-1 flex-col justify-between gap-2 p-3 bg-muted/20">
          <div>
            <p class="text-[10px] font-semibold uppercase tracking-widest text-muted-foreground mb-1">
              Itens
            </p>
            {/* Altura uniforme com scroll sutil para manter todos os cards na mesma altura */}
            <div class="h-16 overflow-y-auto pr-0.5 space-y-1">
              <For each={items()}>
                {(item) => {
                  const meta = ITEM_META[item.item_type] || {
                    icon: Package,
                    label: item.item_type,
                  };
                  const Icon = meta.icon;
                  return (
                    <div class="flex items-center justify-between gap-1 text-[11px]">
                      <span class="flex items-center gap-1 min-w-0 text-foreground">
                        <Icon class="size-3 text-amber-500 shrink-0" />
                        <span class="truncate">{meta.label}</span>
                        <span class="font-mono truncate text-muted-foreground text-[10px]">
                          #{shortItemId(item.item_id)}
                        </span>
                      </span>
                      <span class="font-mono shrink-0 text-muted-foreground text-[10px] text-right">
                        {item.quantity}× {formatMoney(item.unit_price_cents, props.purchase.currency)}
                      </span>
                    </div>
                  );
                }}
              </For>
              <Show when={items().length === 0}>
                <p class="text-[11px] italic text-muted-foreground">
                  Nenhum item
                </p>
              </Show>
            </div>
          </div>

          <div class="space-y-1.5">
            <button
              type="button"
              onClick={() => setOpen((v) => !v)}
              class="flex w-full items-center justify-between rounded-md border border-border bg-card/60 px-2 py-1 text-left transition-colors hover:bg-accent cursor-pointer"
            >
              <span class="flex items-center gap-1.5 text-[11px] text-foreground">
                <Users class="size-3 text-muted-foreground" />
                {attendees().length} {attendees().length === 1 ? "participante" : "participantes"}
              </span>
              <ChevronDown
                class={cn(
                  "size-3.5 text-muted-foreground transition-transform duration-200",
                  open() && "rotate-180",
                )}
              />
            </button>

            <Show when={open()}>
              <ul class="max-h-24 overflow-y-auto space-y-1 pl-1 pr-0.5">
                <For each={attendees()}>
                  {(attendee) => (
                    <li class="flex items-center justify-between gap-2 text-[11px]">
                      <span class="truncate text-foreground">{attendee.name}</span>
                      <span class="font-mono truncate text-muted-foreground text-[10px]">
                        {attendee.email}
                      </span>
                    </li>
                  )}
                </For>
                <Show when={attendees().length === 0}>
                  <li class="text-[11px] italic text-muted-foreground">
                    Nenhum participante
                  </li>
                </Show>
              </ul>
            </Show>
          </div>
        </div>
      </div>
    </Reveal>
  );
}
