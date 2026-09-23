import type { JSX } from "@solidjs/web";
import { Show, createMemo, createSignal } from "solid-js";
import BanIcon from "~icons/lucide/ban";
import CalendarIcon from "~icons/lucide/calendar";
import ClockIcon from "~icons/lucide/clock";
import MailIcon from "~icons/lucide/mail";

import { cn } from "@trieoh/ui-solid";
import { AlertModal } from "@/widgets/ui/AlertModal";
import type { SignatureRequestI, SignatureRequestStatus } from "../model";

const Ban = BanIcon as unknown as (props: { class?: string }) => JSX.Element;
const Calendar = CalendarIcon as unknown as (props: { class?: string }) => JSX.Element;
const Clock = ClockIcon as unknown as (props: { class?: string }) => JSX.Element;
const Mail = MailIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface AdminSignatureRequestCardProps {
  request: SignatureRequestI;
  onCancel?: (reason?: string) => void | Promise<void>;
  isCancelling?: boolean;
}

function StatusPill(props: { status: SignatureRequestStatus }): JSX.Element {
  const meta = createMemo(() => {
    switch (props.status) {
      case "pending":
        return {
          label: "Pendente",
          dotClass: "bg-amber-500 animate-pulse",
          pillClass: "bg-amber-500/10 text-amber-600 dark:text-amber-400 border-amber-500/20",
        };
      case "completed":
        return {
          label: "Concluído",
          dotClass: "bg-emerald-500",
          pillClass: "bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-500/20",
        };
      case "cancelled":
        return {
          label: "Cancelado",
          dotClass: "bg-rose-500",
          pillClass: "bg-rose-500/10 text-rose-600 dark:text-rose-400 border-rose-500/20",
        };
      case "expired":
      default:
        return {
          label: "Expirado",
          dotClass: "bg-muted-foreground",
          pillClass: "bg-muted text-muted-foreground border-border",
        };
    }
  });

  return (
    <span
      class={cn(
        "inline-flex items-center gap-1.5 rounded-full border px-2 py-0.5 text-xs font-medium",
        meta().pillClass,
      )}
    >
      <span class={cn("size-1.5 rounded-full", meta().dotClass)} />
      {meta().label}
    </span>
  );
}

export function AdminSignatureRequestCard(
  props: AdminSignatureRequestCardProps,
): JSX.Element {
  const [cancelOpen, setCancelOpen] = createSignal(false);

  const formattedCreatedDate = () => {
    if (!props.request.created_at) return "";
    try {
      return new Date(props.request.created_at).toLocaleDateString("pt-BR", {
        day: "2-digit",
        month: "short",
        year: "numeric",
      });
    } catch {
      return "";
    }
  };

  const expirationText = createMemo(() => {
    if (!props.request.expires_at) return "";
    try {
      const expiry = new Date(props.request.expires_at).getTime();
      const now = Date.now();
      const diffMs = expiry - now;
      if (diffMs <= 0) return "Expirado";
      const days = Math.ceil(diffMs / (1000 * 60 * 60 * 24));
      return days === 1 ? "Expira hoje" : `Expira em ${days} dias`;
    } catch {
      return "";
    }
  });

  return (
    <>
      <article
        class={cn(
          "group relative flex w-full flex-col justify-between rounded-lg border border-border/60 bg-card p-4 transition-colors duration-150 hover:border-border shadow-xs text-left",
        )}
      >
        <div class="space-y-3">
          {/* Top Row: Status + Actions */}
          <div class="flex items-center justify-between gap-2">
            <StatusPill status={props.request.status} />

            <Show when={props.request.status === "pending" && props.onCancel}>
              <button
                type="button"
                onClick={() => setCancelOpen(true)}
                class="inline-flex items-center gap-1 rounded border border-border/60 bg-background/80 px-2 py-1 text-[11px] font-medium text-muted-foreground transition-colors hover:border-destructive/40 hover:bg-destructive/10 hover:text-destructive cursor-pointer"
                title="Cancelar convite"
              >
                <Ban class="size-3" />
                <span>Cancelar</span>
              </button>
            </Show>
          </div>

          {/* Signatory details */}
          <div class="space-y-1">
            <h3 class="truncate text-base font-semibold text-foreground leading-tight">
              {props.request.signatory_name}
            </h3>

            <Show when={props.request.signatory_title}>
              <p class="truncate text-xs font-medium text-muted-foreground">
                {props.request.signatory_title}
              </p>
            </Show>

            <Show when={props.request.signatory_email}>
              <div class="flex items-center gap-1.5 pt-0.5 text-xs text-muted-foreground truncate">
                <Mail class="size-3.5 shrink-0" />
                <span class="truncate font-mono text-[11px]">
                  {props.request.signatory_email}
                </span>
              </div>
            </Show>
          </div>
        </div>

        {/* Bottom meta row */}
        <div class="mt-4 flex flex-wrap items-center justify-between gap-2 border-t border-border/40 pt-3 text-[11px] text-muted-foreground">
          <div class="flex items-center gap-1">
            <Calendar class="size-3" />
            <span>Enviado em {formattedCreatedDate()}</span>
          </div>

          <Show when={props.request.status === "pending"}>
            <div class="flex items-center gap-1 font-medium text-foreground/80">
              <Clock class="size-3 text-amber-500" />
              <span>{expirationText()}</span>
            </div>
          </Show>
        </div>
      </article>

      {/* Cancel Confirmation Modal */}
      <AlertModal
        open={cancelOpen()}
        onOpenChange={setCancelOpen}
        title="Cancelar convite de assinatura?"
        description={`O convite enviado para ${props.request.signatory_name} (${props.request.signatory_email ?? ""}) será invalidado imediatamente e o link não poderá mais ser preenchido.`}
        confirmLabel="Confirmar cancelamento"
        variant="destructive"
        loading={props.isCancelling}
        onConfirm={async () => {
          await props.onCancel?.("Cancelado pelo administrador");
          setCancelOpen(false);
        }}
      />
    </>
  );
}
