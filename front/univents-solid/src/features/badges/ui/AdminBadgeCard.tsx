import type { JSX } from "@solidjs/web";
import { Show, createSignal } from "solid-js";

import CopyIcon from "~icons/lucide/copy";
import MoreVerticalIcon from "~icons/lucide/more-vertical";
import PencilIcon from "~icons/lucide/pencil";
import Trash2Icon from "~icons/lucide/trash-2";

import { cn } from "@trieoh/ui-solid";
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

  return (
    <>
      <Reveal delay={(props.index ?? 0) * 0.04} animate={props.animate} class="h-full">
        <article
          class={cn(
            "group relative flex h-full w-full min-w-0 flex-col overflow-hidden rounded-lg border border-border/60 bg-card text-left transition-colors duration-150 hover:border-border",
          )}
        >
          {/* Badge Preview Header (Compact Linear style) */}
          <div class="relative flex h-36 w-full items-center justify-center overflow-hidden bg-muted/30 p-3">
            <BadgePreview
              badge={props.item}
              contain
              showVariables={props.kind === "template"}
              ticketName={props.ticketName}
              participantName={props.participantName}
              location={props.location}
              class="max-h-full max-w-full"
            />

            {/* Top-left Pill */}
            <span class="absolute left-2.5 top-2.5 inline-flex items-center rounded border border-border/50 bg-background/80 px-1.5 py-0.5 text-[10px] font-medium text-muted-foreground backdrop-blur-xs">
              {props.kind === "template" ? "Template" : "Crachá"}
            </span>

            {/* Template Actions (Linear subtle icon buttons) */}
            <Show when={props.kind === "template"}>
              <div class="absolute right-2 top-2 flex items-center gap-1">
                <Show when={props.onEdit}>
                  <button
                    type="button"
                    aria-label={`Editar crachá ${title()}`}
                    onClick={() => props.onEdit?.()}
                    class="inline-flex size-6 items-center justify-center rounded border border-border/50 bg-background/80 text-muted-foreground transition-colors hover:border-border hover:text-foreground backdrop-blur-xs cursor-pointer"
                    title="Editar template"
                  >
                    <Pencil class="size-3" />
                  </button>
                </Show>

                <Show when={props.onDuplicate || props.onDelete}>
                  <div class="relative">
                    <button
                      type="button"
                      aria-label={`Ações de ${title()}`}
                      onClick={() => setMenuOpen(!menuOpen())}
                      class="inline-flex size-6 items-center justify-center rounded border border-border/50 bg-background/80 text-muted-foreground transition-colors hover:border-border hover:text-foreground backdrop-blur-xs cursor-pointer"
                    >
                      <MoreVertical class="size-3" />
                    </button>

                    <Show when={menuOpen()}>
                      <div
                        class="fixed inset-0 z-40"
                        onClick={() => setMenuOpen(false)}
                      />
                      <div class="absolute right-0 z-50 mt-1 w-40 rounded-md border border-border bg-popover p-1 text-popover-foreground shadow-sm animate-in fade-in-0 zoom-in-95">
                        <Show when={props.onEdit}>
                          <button
                            type="button"
                            onClick={() => {
                              setMenuOpen(false);
                              props.onEdit?.();
                            }}
                            class="flex w-full items-center gap-2 rounded px-2 py-1.5 text-xs font-medium text-foreground hover:bg-muted cursor-pointer"
                          >
                            <Pencil class="size-3 text-muted-foreground" />
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
                            class="flex w-full items-center gap-2 rounded px-2 py-1.5 text-xs font-medium text-foreground hover:bg-muted cursor-pointer"
                          >
                            <Copy class="size-3 text-muted-foreground" />
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
                            class="flex w-full items-center gap-2 rounded px-2 py-1.5 text-xs font-medium text-destructive hover:bg-destructive/10 cursor-pointer"
                          >
                            <Trash2 class="size-3" />
                            Excluir
                          </button>
                        </Show>
                      </div>
                    </Show>
                  </div>
                </Show>
              </div>
            </Show>
          </div>

          {/* Details Section */}
          <div class="flex flex-1 flex-col justify-between p-3">
            <div class="space-y-0.5">
              <div class="flex items-center justify-between gap-1.5">
                <h3
                  class="truncate text-xs font-medium text-foreground"
                  title={title()}
                >
                  {title()}
                </h3>
                <Show when={isStaff()}>
                  <span class="inline-flex shrink-0 items-center rounded border border-primary/20 bg-primary/10 px-1.5 py-0.5 text-[9px] font-medium text-primary">
                    Staff
                  </span>
                </Show>
              </div>

              <p class="truncate text-[11px] text-muted-foreground" title={subtitle()}>
                {subtitle()}
              </p>
            </div>

            {/* Metadata Footer */}
            <div class="mt-2.5 flex items-center justify-between border-t border-border/40 pt-2 text-[10px] text-muted-foreground">
              <Show
                when={props.kind === "template"}
                fallback={
                  <span>
                    {(props.item as BadgePrintItem).ticket_name || "Ingresso"}
                  </span>
                }
              >
                <span>
                  {badgePxToMm((props.item.design_data as any)?.canvas?.width ?? 321)} ×{" "}
                  {badgePxToMm((props.item.design_data as any)?.canvas?.height ?? 204)} mm
                </span>
                <Show when={props.onEdit}>
                  <button
                    type="button"
                    onClick={() => props.onEdit?.()}
                    class="text-[11px] font-medium text-foreground hover:text-primary transition-colors cursor-pointer"
                  >
                    Editar
                  </button>
                </Show>
              </Show>
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
