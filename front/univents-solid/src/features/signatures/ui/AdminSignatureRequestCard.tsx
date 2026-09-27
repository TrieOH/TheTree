import type { JSX } from "@solidjs/web";
import { Button, cn } from "@trieoh/ui-solid";
import { Show, createMemo, createSignal } from "solid-js";
import CalendarIcon from "~icons/lucide/calendar";
import ClockIcon from "~icons/lucide/clock";
import MailIcon from "~icons/lucide/mail";

import { AlertModal } from "@/widgets/ui/AlertModal";
import type { SignatureRequestI, SignatureRequestStatus } from "../model";

const Calendar = CalendarIcon as unknown as (props: { class?: string }) => JSX.Element;
const Clock = ClockIcon as unknown as (props: { class?: string }) => JSX.Element;
const Mail = MailIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface AdminSignatureRequestCardProps {
  request: SignatureRequestI;
  onCancel?: (reason?: string) => Promise<void> | void;
  isCancelling?: boolean;
}

function StatusBadge(props: { status: SignatureRequestStatus }): JSX.Element {
  const meta = createMemo(() => {
    switch (props.status) {
      case "completed":
        return {
          label: "Concluída",
          pillClass: "border-emerald-500/20 bg-emerald-500/10 text-emerald-600 dark:text-emerald-400",
          dotClass: "bg-emerald-500",
        };
      case "expired":
        return {
          label: "Expirada",
          pillClass: "border-amber-500/20 bg-amber-500/10 text-amber-600 dark:text-amber-400",
          dotClass: "bg-amber-500",
        };
      case "cancelled":
        return {
          label: "Cancelada",
          pillClass: "border-destructive/20 bg-destructive/10 text-destructive",
          dotClass: "bg-destructive",
        };
      case "pending":
      default:
        return {
          label: "Pendente",
          pillClass: "border-primary/20 bg-primary/10 text-primary",
          dotClass: "bg-primary animate-pulse",
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

      if (diffMs <= 0) return "Prazo expirado";
      const days = Math.ceil(diffMs / (1000 * 60 * 60 * 24));
      if (days === 1) return "Expira hoje";
      return `Expira em ${days} dias`;
    } catch {
      return "";
    }
  });

  const handleConfirmCancel = async () => {
    await props.onCancel?.("Cancelado pelo administrador");
    setCancelOpen(false);
  };

  return (
    <>
      <article class="flex flex-col justify-between rounded-xl border border-border bg-card p-4 transition-all duration-200 hover:border-foreground/20 hover:shadow-xs text-left">
        <div>
          {/* Header Row */}
          <div class="flex items-start justify-between gap-3">
            <div class="min-w-0 space-y-0.5">
              <h4
                class="font-semibold text-sm text-foreground line-clamp-1"
                title={props.request.signatory_name}
              >
                {props.request.signatory_name}
              </h4>
              <Show when={props.request.signatory_title}>
                <p
                  class="text-xs text-muted-foreground line-clamp-1"
                  title={props.request.signatory_title!}
                >
                  {props.request.signatory_title}
                </p>
              </Show>
            </div>

            <StatusBadge status={props.request.status} />
          </div>

          {/* Details */}
          <div class="mt-3.5 space-y-1 text-[11px] text-muted-foreground">
            <Show when={props.request.signatory_email}>
              <div class="flex items-center gap-1.5 truncate">
                <Mail class="size-3 shrink-0" />
                <span class="truncate">{props.request.signatory_email}</span>
              </div>
            </Show>

            <Show when={formattedCreatedDate()}>
              <div class="flex items-center gap-1.5">
                <Calendar class="size-3 shrink-0" />
                <span>Enviado em {formattedCreatedDate()}</span>
              </div>
            </Show>

            <Show when={props.request.status === "pending" && expirationText()}>
              <div class="flex items-center gap-1.5 text-amber-600 dark:text-amber-400">
                <Clock class="size-3 shrink-0" />
                <span>{expirationText()}</span>
              </div>
            </Show>

            <Show when={props.request.status_reason}>
              <p class="mt-2 text-xs italic text-muted-foreground bg-muted/40 p-2 rounded-md">
                "{props.request.status_reason}"
              </p>
            </Show>
          </div>
        </div>

        {/* Footer Actions */}
        <div class="mt-4 flex items-center justify-between border-t border-border/60 pt-3">
          <Show
            when={props.request.status === "pending"}
            fallback={
              <span class="text-[11px] text-muted-foreground">
                Solicitação encerrada
              </span>
            }
          >
            <Button
              variant="outline"
              size="sm"
              onClick={() => setCancelOpen(true)}
              class="h-7 text-xs text-destructive hover:bg-destructive/10 hover:text-destructive hover:border-destructive/30 cursor-pointer"
            >
              Cancelar solicitação
            </Button>
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
        onConfirm={handleConfirmCancel}
      />
    </>
  );
}
