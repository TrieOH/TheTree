import type { JSX } from "@solidjs/web";
import { Show, createSignal } from "solid-js";

import CopyIcon from "~icons/lucide/copy";
import MoreVerticalIcon from "~icons/lucide/more-vertical";
import PencilIcon from "~icons/lucide/pencil";
import Trash2Icon from "~icons/lucide/trash-2";

import { Button, cn } from "@trieoh/ui-solid";
import { Reveal } from "@/shared/ui/Reveal";
import { AlertModal } from "@/widgets/ui/AlertModal";
import { type BadgePrintItem, type BadgeTemplate, badgePxToMm } from "../model";
import { BadgePreview } from "./BadgePreview";

const Copy = CopyIcon as unknown as (props: { class?: string }) => JSX.Element;
const MoreVertical = MoreVerticalIcon as unknown as (props: { class?: string }) => JSX.Element;
const Pencil = PencilIcon as unknown as (props: { class?: string }) => JSX.Element;
const Trash2 = Trash2Icon as unknown as (props: { class?: string }) => JSX.Element;

export interface AdminBadgeCardProps {
  item: BadgeTemplate | BadgePrintItem;
  kind: "template" | "emission";
  index?: number;
  animate?: boolean;
  ticketName?: string;
  participantName?: string;
  location?: string;
  onEdit?: () => void;
  onDelete?: () => void;
  onDuplicate?: () => void;
  onView?: () => void;
}

export function AdminBadgeCard(props: AdminBadgeCardProps): JSX.Element {
  const [deleteOpen, setDeleteOpen] = createSignal(false);
  const [menuOpen, setMenuOpen] = createSignal(false);

  const title = () =>
    props.kind === "template"
      ? (props.item as BadgeTemplate).name
      : (props.item as BadgePrintItem).event_name;

  const subtitle = () =>
    props.kind === "template"
      ? props.ticketName ?? "Padrão da edição"
      : (props.item as BadgePrintItem).edition_name;

  const isStaff = () =>
    props.kind === "template" && (props.item as BadgeTemplate).origin === "staff";

  const handleCardClick = () => {
    if (props.onView) {
      props.onView();
    } else if (props.onEdit) {
      props.onEdit();
    }
  };

  return (
    <>
      <Reveal delay={(props.index ?? 0) * 0.04} animate={props.animate} class="h-full">
        <article
          role="button"
          tabindex={0}
          aria-label={`Visualizar crachá ${title()}`}
          onClick={handleCardClick}
          onKeyDown={(e) => {
            if (e.key === "Enter" || e.key === " ") {
              e.preventDefault();
              handleCardClick();
            }
          }}
          class={cn(
            "group relative flex h-full w-full min-w-0 flex-col overflow-hidden rounded-xl border border-border bg-card text-left shadow-xs transition-all duration-200",
            "hover:-translate-y-0.5 hover:border-foreground/20 hover:shadow-md cursor-pointer",
            "focus:outline-none focus-visible:ring-2 focus-visible:ring-ring",
          )}
        >
          {/* Badge Preview Header */}
          <div class="relative flex h-64 w-full items-center justify-center overflow-hidden bg-muted p-4">
            <BadgePreview
              badge={props.item}
              contain
              showVariables={props.kind === "template"}
              ticketName={props.ticketName}
              participantName={props.participantName}
              location={props.location}
              class="max-h-full max-w-full"
            />

            {/* Template Actions Menu */}
            <Show when={props.kind === "template"}>
              <div class="absolute right-2.5 top-2.5" onClick={(e) => e.stopPropagation()}>
                <div class="relative">
                  <Button
                    type="button"
                    variant="outline"
                    size="icon"
                    aria-label={`Ações de ${title()}`}
                    onClick={() => setMenuOpen(!menuOpen())}
                    class="size-8 rounded-full bg-background/90 text-foreground shadow-xs backdrop-blur-xs hover:bg-background"
                  >
                    <MoreVertical class="size-4" />
                  </Button>

                  <Show when={menuOpen()}>
                    <div
                      class="fixed inset-0 z-40"
                      onClick={() => setMenuOpen(false)}
                    />
                    <div class="absolute right-0 z-50 mt-1 w-44 rounded-lg border border-border bg-popover p-1 text-popover-foreground shadow-md animate-in fade-in-0 zoom-in-95">
                      <Show when={props.onEdit}>
                        <button
                          type="button"
                          onClick={() => {
                            setMenuOpen(false);
                            props.onEdit?.();
                          }}
                          class="flex w-full items-center gap-2 rounded-md px-2.5 py-1.5 text-xs font-medium text-foreground hover:bg-muted"
                        >
                          <Pencil class="size-3.5 text-muted-foreground" />
                          Editar
                        </button>
                      </Show>
                      <Show when={props.onDuplicate}>
                        <button
                          type="button"
                          onClick={() => {
                            setMenuOpen(false);
                            props.onDuplicate?.();
                          }}
                          class="flex w-full items-center gap-2 rounded-md px-2.5 py-1.5 text-xs font-medium text-foreground hover:bg-muted"
                        >
                          <Copy class="size-3.5 text-muted-foreground" />
                          Duplicar
                        </button>
                      </Show>
                      <Show when={props.onDelete}>
                        <div class="my-1 border-t border-border" />
                        <button
                          type="button"
                          onClick={() => {
                            setMenuOpen(false);
                            setDeleteOpen(true);
                          }}
                          class="flex w-full items-center gap-2 rounded-md px-2.5 py-1.5 text-xs font-medium text-destructive hover:bg-destructive/10"
                        >
                          <Trash2 class="size-3.5" />
                          Excluir
                        </button>
                      </Show>
                    </div>
                  </Show>
                </div>
              </div>
            </Show>

            <span class="absolute left-3 top-3 inline-flex items-center gap-2 rounded-full border border-border/70 bg-background/85 px-3 py-1 text-[10px] font-semibold text-muted-foreground shadow-xs">
              {props.kind === "template" ? "Template" : "Crachá"}
            </span>
          </div>

          {/* Details Section */}
          <div class="flex flex-1 flex-col justify-between p-4">
            <div class="space-y-1.5">
              <div class="flex items-start justify-between gap-2">
                <h3
                  class="truncate text-sm font-semibold tracking-tight text-foreground transition-colors group-hover:text-primary"
                  title={title()}
                >
                  {title()}
                </h3>
                <Show when={isStaff()}>
                  <span class="inline-flex shrink-0 items-center rounded-full bg-primary/10 px-2 py-0.5 text-[10px] font-semibold text-primary">
                    Staff
                  </span>
                </Show>
              </div>

              <p class="truncate text-xs text-muted-foreground" title={subtitle()}>
                {subtitle()}
              </p>
            </div>

            {/* Template Dimensions Footer */}
            <div class="mt-4 flex items-center justify-between border-t border-border/40 pt-3 text-[11px] text-muted-foreground">
              <span>
                {badgePxToMm((props.item.design_data as any)?.canvas?.width ?? 321)} × {badgePxToMm((props.item.design_data as any)?.canvas?.height ?? 204)} mm
              </span>
              <span class="text-primary hover:underline">
                {props.kind === "template" ? "Configurar →" : "Visualizar →"}
              </span>
            </div>
          </div>
        </article>
      </Reveal>

      {/* Delete Confirmation Dialog */}
      <AlertModal
        open={deleteOpen()}
        onOpenChange={setDeleteOpen}
        title="Excluir crachá"
        description={`Tem certeza que deseja excluir o template "${title()}"? Esta ação não pode ser desfeita.`}
        confirmLabel="Excluir"
        variant="destructive"
        onConfirm={() => {
          setDeleteOpen(false);
          props.onDelete?.();
        }}
      />
    </>
  );
}
