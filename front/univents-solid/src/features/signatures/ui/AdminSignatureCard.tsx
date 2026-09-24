import type { JSX } from "@solidjs/web";
import { Show, createSignal } from "solid-js";
import CalendarIcon from "~icons/lucide/calendar";
import MailIcon from "~icons/lucide/mail";
import Trash2Icon from "~icons/lucide/trash-2";

import { AlertModal } from "@/widgets/ui/AlertModal";
import type { SignatureI } from "../model";

const Calendar = CalendarIcon as unknown as (props: { class?: string }) => JSX.Element;
const Mail = MailIcon as unknown as (props: { class?: string }) => JSX.Element;
const Trash2 = Trash2Icon as unknown as (props: { class?: string }) => JSX.Element;

export interface AdminSignatureCardProps {
  signature: SignatureI;
  onDelete?: () => Promise<void> | void;
  isDeleting?: boolean;
}

export function AdminSignatureCard(
  props: AdminSignatureCardProps,
): JSX.Element {
  const [deleteOpen, setDeleteOpen] = createSignal(false);

  const formattedDate = () => {
    try {
      return new Date(props.signature.created_at).toLocaleDateString("pt-BR", {
        day: "2-digit",
        month: "short",
        year: "numeric",
      });
    } catch {
      return "";
    }
  };

  const handleConfirmDelete = async () => {
    await props.onDelete?.();
    setDeleteOpen(false);
  };

  return (
    <>
      <article class="group relative flex flex-col overflow-hidden rounded-xl border border-border bg-card transition-all duration-200 hover:border-foreground/20 hover:shadow-xs">
        {/* Preview Frame with transparency pattern */}
        <div class="relative flex h-40 w-full items-center justify-center overflow-hidden border-b border-border/80 bg-white p-6">
          <div
            class="pointer-events-none absolute inset-0 opacity-40"
            style={{
              "background-image":
                "radial-gradient(#94a3b8 0.75px, transparent 0.75px), radial-gradient(#94a3b8 0.75px, #ffffff 0.75px)",
              "background-size": "16px 16px",
              "background-position": "0 0, 8px 8px",
            }}
          />

          <img
            src={props.signature.image_url}
            alt={`Assinatura de ${props.signature.signatory_name}`}
            class="max-h-full max-w-full object-contain filter drop-shadow-xs transition-transform duration-200 group-hover:scale-[1.02]"
            loading="lazy"
          />

          {/* Top-left Pill */}
          <span class="absolute left-2.5 top-2.5 z-10 inline-flex items-center rounded border border-border/60 bg-background/90 px-2 py-0.5 text-[10px] font-medium text-muted-foreground shadow-xs backdrop-blur-xs">
            Assinatura
          </span>

          {/* Top-right Actions */}
          <Show when={props.onDelete}>
            <div class="absolute right-2.5 top-2.5 z-10">
              <button
                type="button"
                aria-label={`Excluir assinatura de ${props.signature.signatory_name}`}
                onClick={() => setDeleteOpen(true)}
                class="flex size-7 items-center justify-center rounded-md border border-border/60 bg-background/90 text-muted-foreground shadow-xs backdrop-blur-xs transition-colors hover:border-destructive/40 hover:bg-destructive/10 hover:text-destructive cursor-pointer"
                title="Excluir assinatura"
              >
                <Trash2 class="size-3.5" />
              </button>
            </div>
          </Show>
        </div>

        {/* Content Body */}
        <div class="flex flex-1 flex-col justify-between p-4 text-left">
          <div class="space-y-1">
            <h4
              class="font-semibold text-sm text-foreground line-clamp-1"
              title={props.signature.signatory_name}
            >
              {props.signature.signatory_name}
            </h4>
            <Show when={props.signature.signatory_title}>
              <p
                class="text-xs text-muted-foreground line-clamp-1"
                title={props.signature.signatory_title!}
              >
                {props.signature.signatory_title}
              </p>
            </Show>
          </div>

          <div class="mt-4 flex flex-col gap-1 border-t border-border/60 pt-3 text-[11px] text-muted-foreground">
            <Show when={props.signature.signatory_email}>
              <div class="flex items-center gap-1.5 truncate">
                <Mail class="size-3 shrink-0" />
                <span class="truncate">{props.signature.signatory_email}</span>
              </div>
            </Show>

            <Show when={formattedDate()}>
              <div class="flex items-center gap-1.5">
                <Calendar class="size-3 shrink-0" />
                <span>Adicionada em {formattedDate()}</span>
              </div>
            </Show>
          </div>
        </div>
      </article>

      {/* Delete Confirmation Modal */}
      <AlertModal
        open={deleteOpen()}
        onOpenChange={setDeleteOpen}
        title="Remover assinatura?"
        description={`A assinatura de "${props.signature.signatory_name}" será removida desta edição. Certificados já emitidos que utilizam a imagem continuarão com o arquivo preservado.`}
        confirmLabel="Remover assinatura"
        variant="destructive"
        loading={props.isDeleting}
        onConfirm={handleConfirmDelete}
      />
    </>
  );
}
