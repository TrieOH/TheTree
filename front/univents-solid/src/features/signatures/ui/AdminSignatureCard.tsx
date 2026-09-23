import type { JSX } from "@solidjs/web";
import { Show, createSignal } from "solid-js";
import CopyIcon from "~icons/lucide/copy";
import CheckIcon from "~icons/lucide/check";
import Trash2Icon from "~icons/lucide/trash-2";
import MailIcon from "~icons/lucide/mail";

import { cn } from "@trieoh/ui-solid";
import { AlertModal } from "@/widgets/ui/AlertModal";
import type { SignatureI } from "../model";

const Copy = CopyIcon as unknown as (props: { class?: string }) => JSX.Element;
const Check = CheckIcon as unknown as (props: { class?: string }) => JSX.Element;
const Trash2 = Trash2Icon as unknown as (props: { class?: string }) => JSX.Element;
const Mail = MailIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface AdminSignatureCardProps {
  signature: SignatureI;
  onDelete?: () => void | Promise<void>;
  isDeleting?: boolean;
}

export function AdminSignatureCard(props: AdminSignatureCardProps): JSX.Element {
  const [deleteOpen, setDeleteOpen] = createSignal(false);
  const [copied, setCopied] = createSignal(false);

  const formattedDate = () => {
    if (!props.signature.created_at) return "";
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

  const handleCopyUrl = async () => {
    try {
      await navigator.clipboard.writeText(props.signature.image_url);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch {
      // ignore
    }
  };

  return (
    <>
      <article
        class={cn(
          "group relative flex h-full w-full min-w-0 flex-col rounded-lg border border-border/60 bg-card text-left transition-colors duration-150 hover:border-border shadow-xs",
        )}
      >
        {/* Preview Container */}
        <div class="relative flex h-36 w-full items-center justify-center overflow-hidden rounded-t-lg bg-white p-4">
          {/* Subtle checkered backdrop pattern for transparency awareness */}
          <div
            class="absolute inset-0 opacity-[0.03] pointer-events-none"
            style={{
              "background-image":
                "radial-gradient(#000 1px, transparent 1px), radial-gradient(#000 1px, transparent 1px)",
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
          <span class="absolute left-2.5 top-2.5 inline-flex items-center rounded border border-border/50 bg-background/80 px-2 py-0.5 text-[10px] font-medium text-muted-foreground backdrop-blur-xs">
            Assinatura
          </span>

          {/* Top-right Actions */}
          <div class="absolute right-2.5 top-2.5 z-10 flex items-center gap-1 opacity-80 group-hover:opacity-100 transition-opacity">
            <button
              type="button"
              aria-label="Copiar link da assinatura"
              onClick={handleCopyUrl}
              class="inline-flex size-7 items-center justify-center rounded border border-border/50 bg-background/80 text-muted-foreground transition-colors hover:border-border hover:text-foreground backdrop-blur-xs cursor-pointer"
              title="Copiar URL da imagem"
            >
              <Show when={copied()} fallback={<Copy class="size-3.5" />}>
                <Check class="size-3.5 text-emerald-600 dark:text-emerald-400" />
              </Show>
            </button>

            <Show when={props.onDelete}>
              <button
                type="button"
                aria-label={`Excluir assinatura de ${props.signature.signatory_name}`}
                onClick={() => setDeleteOpen(true)}
                class="inline-flex size-7 items-center justify-center rounded border border-destructive/30 bg-background/80 text-destructive/80 transition-colors hover:border-destructive hover:bg-destructive/10 hover:text-destructive backdrop-blur-xs cursor-pointer"
                title="Excluir assinatura"
              >
                <Trash2 class="size-3.5" />
              </button>
            </Show>
          </div>
        </div>

        {/* Info Content */}
        <div class="flex flex-1 flex-col justify-between p-3.5">
          <div class="space-y-1">
            <h3 class="truncate text-sm font-semibold text-foreground leading-tight">
              {props.signature.signatory_name}
            </h3>

            <Show when={props.signature.signatory_title}>
              <p class="truncate text-xs font-medium text-muted-foreground">
                {props.signature.signatory_title}
              </p>
            </Show>

            <Show when={props.signature.signatory_email}>
              <div class="flex items-center gap-1.5 pt-0.5 text-xs text-muted-foreground/80 truncate">
                <Mail class="size-3 shrink-0" />
                <span class="truncate">{props.signature.signatory_email}</span>
              </div>
            </Show>
          </div>

          <div class="mt-3 flex items-center justify-between border-t border-border/40 pt-2 text-[11px] text-muted-foreground">
            <span>Criado em</span>
            <span class="font-mono">{formattedDate()}</span>
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
        onConfirm={async () => {
          await props.onDelete?.();
          setDeleteOpen(false);
        }}
      />
    </>
  );
}
