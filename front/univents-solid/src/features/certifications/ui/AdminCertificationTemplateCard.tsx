import type { JSX } from "@solidjs/web";
import { Show, createSignal } from "solid-js";
import CopyIcon from "~icons/lucide/copy";
import EyeIcon from "~icons/lucide/eye";
import MoreVerticalIcon from "~icons/lucide/more-vertical";
import PencilIcon from "~icons/lucide/pencil";
import Trash2Icon from "~icons/lucide/trash-2";

import { cn } from "@trieoh/ui-solid";
import { AlertModal } from "@/widgets/ui/AlertModal";
import { Reveal } from "@/shared/ui/Reveal";
import type { CertificationTemplateI } from "../model";
import { CertificatePreview } from "./CertificatePreview";

const Copy = CopyIcon as unknown as (props: { class?: string }) => JSX.Element;
const Eye = EyeIcon as unknown as (props: { class?: string }) => JSX.Element;
const MoreVertical = MoreVerticalIcon as unknown as (props: { class?: string }) => JSX.Element;
const Pencil = PencilIcon as unknown as (props: { class?: string }) => JSX.Element;
const Trash2 = Trash2Icon as unknown as (props: { class?: string }) => JSX.Element;

export interface AdminCertificationTemplateCardProps {
  template: CertificationTemplateI;
  index?: number;
  animate?: boolean;
  onEdit?: () => void;
  onView?: () => void;
  onDuplicate?: () => void;
  onDelete?: () => void;
}

export function AdminCertificationTemplateCard(
  props: AdminCertificationTemplateCardProps,
): JSX.Element {
  const [deleteOpen, setDeleteOpen] = createSignal(false);
  const [menuOpen, setMenuOpen] = createSignal(false);

  const title = () => props.template.name;
  const subtitle = () =>
    props.template.description ||
    (props.template.design_data?.background
      ? "Com imagem de fundo"
      : "Certificado padrão");

  const canvas = () => props.template.design_data?.canvas;

  return (
    <>
      <Reveal
        delay={(props.index ?? 0) * 0.04}
        animate={props.animate}
        class="h-full"
      >
        <article
          class={cn(
            "group relative flex h-full w-full min-w-0 flex-col rounded-lg border border-border/60 bg-card text-left transition-colors duration-150 hover:border-border shadow-xs overflow-hidden",
          )}
        >
          {/* Top row: certificate preview & actions */}
          <div class="relative flex h-36 sm:h-40 w-full items-center justify-center overflow-hidden rounded-t-lg bg-muted/20 p-2.5 sm:p-3 border-b border-border/40">
            <span class="absolute left-2 top-2 z-10 inline-flex items-center rounded border border-border/60 bg-background/90 px-1.5 py-0.5 text-[10px] font-medium text-muted-foreground shadow-xs backdrop-blur-xs">
              Template
            </span>

            <CertificatePreview
              template={props.template}
              contain
              framed={false}
              showVariables={true}
              class="max-h-full max-w-full shadow-sm transition-transform duration-200 group-hover:scale-[1.02]"
            />

            {/* Actions Dropdown for templates */}
            <div class="absolute top-2 right-2 z-10">
              <div class="relative">
                <button
                  type="button"
                  onClick={() => setMenuOpen(!menuOpen())}
                  class="flex size-7 items-center justify-center rounded-md border border-border/60 bg-background/90 text-muted-foreground shadow-xs backdrop-blur-xs transition-colors hover:border-border hover:bg-background hover:text-foreground cursor-pointer"
                  aria-label={`Opções do modelo ${title()}`}
                >
                  <MoreVertical class="size-4" />
                </button>

                <Show when={menuOpen()}>
                  <div
                    class="absolute right-0 top-8 z-30 min-w-36 rounded-md border border-border bg-popover p-1 shadow-md"
                    onFocusOut={(e) => {
                      if (!e.currentTarget.contains(e.relatedTarget as Node)) {
                        setMenuOpen(false);
                      }
                    }}
                  >
                    <Show when={props.onView}>
                      <button
                        type="button"
                        onClick={() => {
                          setMenuOpen(false);
                          props.onView?.();
                        }}
                        class="flex w-full items-center gap-2 rounded px-2 py-1.5 text-xs text-foreground hover:bg-accent cursor-pointer transition-colors"
                      >
                        <Eye class="size-3.5 text-muted-foreground" />
                        <span>Visualizar</span>
                      </button>
                    </Show>
                    <Show when={props.onEdit}>
                      <button
                        type="button"
                        onClick={() => {
                          setMenuOpen(false);
                          props.onEdit?.();
                        }}
                        class="flex w-full items-center gap-2 rounded px-2 py-1.5 text-xs text-foreground hover:bg-accent cursor-pointer transition-colors"
                      >
                        <Pencil class="size-3.5 text-muted-foreground" />
                        <span>Editar</span>
                      </button>
                    </Show>
                    <Show when={props.onDuplicate}>
                      <button
                        type="button"
                        onClick={() => {
                          setMenuOpen(false);
                          props.onDuplicate?.();
                        }}
                        class="flex w-full items-center gap-2 rounded px-2 py-1.5 text-xs text-foreground hover:bg-accent cursor-pointer transition-colors"
                      >
                        <Copy class="size-3.5 text-muted-foreground" />
                        <span>Duplicar</span>
                      </button>
                    </Show>
                    <Show when={props.onDelete}>
                      <button
                        type="button"
                        onClick={() => {
                          setMenuOpen(false);
                          setDeleteOpen(true);
                        }}
                        class="flex w-full items-center gap-2 rounded px-2 py-1.5 text-xs text-destructive hover:bg-destructive/10 cursor-pointer transition-colors"
                      >
                        <Trash2 class="size-3.5" />
                        <span>Excluir</span>
                      </button>
                    </Show>
                  </div>
                </Show>
              </div>
            </div>
          </div>

          {/* Info Content */}
          <div class="flex flex-1 flex-col justify-between p-3">
            <div class="space-y-0.5">
              <div class="flex items-center justify-between gap-2">
                <h4
                  class="truncate text-xs font-semibold text-foreground leading-tight"
                  title={title()}
                >
                  {title()}
                </h4>
              </div>
              <p
                class="truncate text-[11px] text-muted-foreground"
                title={subtitle()}
              >
                {subtitle()}
              </p>
            </div>

            {/* Footer */}
            <div class="mt-2.5 flex items-center justify-between border-t border-border/40 pt-2 text-[11px] text-muted-foreground">
              <span>
                {canvas()?.width ?? 1484} × {canvas()?.height ?? 1060} px
              </span>
              <Show when={props.onEdit}>
                <button
                  type="button"
                  onClick={() => props.onEdit?.()}
                  aria-label={`Editar template ${title()}`}
                  class="text-[11px] font-medium text-foreground hover:text-primary transition-colors cursor-pointer"
                >
                  Editar
                </button>
              </Show>
            </div>
          </div>
        </article>
      </Reveal>

      {/* Delete Confirmation Modal */}
      <AlertModal
        open={deleteOpen()}
        onOpenChange={setDeleteOpen}
        title="Excluir modelo de certificado"
        description={`Tem certeza que deseja remover o modelo "${title()}"? Esta ação não pode ser desfeita.`}
        confirmLabel="Excluir modelo"
        variant="destructive"
        onConfirm={() => {
          setDeleteOpen(false);
          props.onDelete?.();
        }}
      />
    </>
  );
}
