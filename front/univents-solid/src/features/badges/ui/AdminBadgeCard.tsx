import type { JSX } from "@solidjs/web";
import { Show, createSignal } from "solid-js";

import CopyIcon from "~icons/lucide/copy";
import MoreVerticalIcon from "~icons/lucide/more-vertical";
import PencilIcon from "~icons/lucide/pencil";
import Trash2Icon from "~icons/lucide/trash-2";

import { cn } from "@trieoh/ui-solid";
import { Reveal } from "@/shared/ui/Reveal";
import { AlertModal } from "@/widgets/ui/AlertModal";
import {
  type BadgeDesign,
  type BadgePrintItem,
  type BadgeTemplate,
  badgePxToMm,
} from "../model";
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
  onPrint?: () => void;
  onView?: () => void;
}

export function AdminBadgeCard(props: AdminBadgeCardProps): JSX.Element {
  const [showMenu, setShowMenu] = createSignal(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = createSignal(false);

  const title = () =>
    props.kind === "template"
      ? (props.item as BadgeTemplate).name
      : (props.item as BadgePrintItem).event_name || "Participante";

  const subtitle = () =>
    props.kind === "template"
      ? props.ticketName || "Geral (todos os ingressos)"
      : (props.item as BadgePrintItem).ticket_name || "Ingresso";

  return (
    <>
      <Reveal
        delay={props.animate ? (props.index ?? 0) * 0.04 : 0}
        animate={props.animate}
        class="h-full"
      >
        <article
          class={cn(
            "group relative flex h-full flex-col justify-between overflow-hidden rounded-lg border border-border/60 bg-card p-3.5 text-left transition-colors duration-150 hover:border-border",
          )}
        >
          <div class="space-y-3">
            {/* Top row: badge preview & actions */}
            <div class="relative flex items-center justify-center rounded-md bg-muted/40 p-4 border border-border/40">
              <span class="absolute left-2.5 top-2.5 inline-flex items-center rounded border border-border/70 bg-background/85 px-2 py-0.5 text-[10px] font-medium text-muted-foreground shadow-xs">
                {props.kind === "template" ? "Template" : "Crachá"}
              </span>
              <BadgePreview
                badge={props.item}
                ticketName={props.ticketName}
                participantName={props.participantName}
                location={props.location}
                framed={false}
                class="w-36 shadow-sm"
              />

              {/* Actions Dropdown for templates */}
              <Show when={props.kind === "template"}>
                <div class="absolute top-2 right-2">
                  <div class="relative">
                    <button
                      type="button"
                      onClick={() => setShowMenu(!showMenu())}
                      class="flex size-7 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-background hover:text-foreground cursor-pointer"
                      aria-label="Opções do modelo"
                    >
                      <MoreVertical class="size-4" />
                    </button>

                    <Show when={showMenu()}>
                      <div
                        class="absolute right-0 top-8 z-30 min-w-36 rounded-md border border-border bg-popover p-1 shadow-md"
                        onFocusOut={(e) => {
                          if (!e.currentTarget.contains(e.relatedTarget as Node)) {
                            setShowMenu(false);
                          }
                        }}
                      >
                        <Show when={props.onEdit}>
                          <button
                            type="button"
                            onClick={() => {
                              setShowMenu(false);
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
                              setShowMenu(false);
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
                              setShowMenu(false);
                              setShowDeleteConfirm(true);
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
              </Show>
            </div>

            {/* Info */}
            <div class="space-y-1">
              <div class="flex items-center justify-between gap-2">
                <h4
                  class="truncate text-xs font-medium text-foreground"
                  title={title()}
                >
                  {title()}
                </h4>
                <Show when={props.kind === "template" && (props.item as BadgeTemplate).origin}>
                  <span class="rounded bg-primary/10 px-1.5 py-0.5 text-[10px] font-medium text-primary">
                    {(props.item as BadgeTemplate).origin}
                  </span>
                </Show>
              </div>
              <p
                class="truncate text-[11px] text-muted-foreground"
                title={subtitle()}
              >
                {subtitle()}
              </p>
            </div>
          </div>

          {/* Footer */}
          <div class="mt-4 pt-2.5 border-t border-border/40">
            <div class="flex items-center justify-between text-[11px] text-muted-foreground">
              <Show
                when={props.kind === "template"}
                fallback={
                  <span>
                    {(props.item as BadgePrintItem).ticket_name || "Ingresso"}
                  </span>
                }
              >
                <span>
                  {badgePxToMm((props.item.design_data as BadgeDesign | undefined)?.canvas?.width ?? 321)} ×{" "}
                  {badgePxToMm((props.item.design_data as BadgeDesign | undefined)?.canvas?.height ?? 204)} mm
                </span>
                <Show when={props.onEdit}>
                  <button
                    type="button"
                    onClick={() => props.onEdit?.()}
                    aria-label={`Editar crachá ${title()}`}
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

      {/* Delete Confirmation Modal */}
      <AlertModal
        open={showDeleteConfirm()}
        onOpenChange={setShowDeleteConfirm}
        title="Excluir modelo de crachá"
        description={`Tem certeza que deseja remover o modelo "${title()}"? Esta ação não pode ser desfeita.`}
        confirmLabel="Excluir modelo"
        variant="destructive"
        onConfirm={() => {
          setShowDeleteConfirm(false);
          props.onDelete?.();
        }}
      />
    </>
  );
}
