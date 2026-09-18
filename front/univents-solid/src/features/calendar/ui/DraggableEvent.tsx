import type { JSX } from "@solidjs/web";
import { Show, createEffect, createSignal, onCleanup } from "solid-js";

import GiftIcon from "~icons/lucide/gift";
import PencilIcon from "~icons/lucide/pencil";
import Trash2Icon from "~icons/lucide/trash-2";
import UserCheckIcon from "~icons/lucide/user-check";

import { formatTimeRange, getDurationMinutes } from "../lib/date";
import {
  endCalendarDrag,
  shouldIgnoreClick,
  startCalendarDrag,
} from "../lib/drag-state";
import type { EventColor, OccurrenceI, ProgramI } from "../model";

const Gift = GiftIcon as unknown as (props: { class?: string }) => JSX.Element;
const Pencil = PencilIcon as unknown as (props: { class?: string }) => JSX.Element;
const Trash2 = Trash2Icon as unknown as (props: { class?: string }) => JSX.Element;
const UserCheck = UserCheckIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface DraggableEventProps {
  occurrence: OccurrenceI;
  program: ProgramI | undefined;
  color: EventColor;
  onClick: (occurrence: OccurrenceI) => void;
  onDelete?: (occurrenceId: string) => void;
  onOpenAttendance?: (occurrence: OccurrenceI) => void;
  onOpenDraw?: (occurrence: OccurrenceI) => void;
  overlapIndex?: number;
  overlapCount?: number;
  isContinuation?: boolean;
  continuesNextDay?: boolean;
}

export function DraggableEvent(props: DraggableEventProps): JSX.Element {
  const [menuOpen, setMenuOpen] = createSignal(false);
  let menuRef: HTMLDivElement | null = null;

  const overlapIndex = () => props.overlapIndex ?? 0;
  const overlapCount = () => props.overlapCount ?? 1;

  const start = () => new Date(props.occurrence.starts_at);
  const startMins = () => start().getHours() * 60 + start().getMinutes();
  const durMins = () =>
    Math.max(
      15,
      getDurationMinutes(props.occurrence.starts_at, props.occurrence.ends_at),
    );

  // Absolute positioning in the 1440px day timeline (1 min = 1px)
  const top = () => startMins();
  const height = () =>
    Math.min(
      Math.max(22, durMins() - 2),
      24 * 60 - startMins(),
    );

  createEffect(
    () => menuOpen(),
    (open) => {
      if (!open) return;

      const handleClickOutside = (e: MouseEvent) => {
        if (menuRef && !menuRef.contains(e.target as Node)) {
          setMenuOpen(false);
        }
      };

      document.addEventListener("mousedown", handleClickOutside);
      onCleanup(() => {
        document.removeEventListener("mousedown", handleClickOutside);
      });
    },
  );

  return (
    <div
      draggable="true"
      onDragStart={(e) => {
        e.stopPropagation();
        startCalendarDrag(e, {
          type: "occurrence",
          occurrenceId: props.occurrence.id,
        });
      }}
      onDragEnd={() => {
        endCalendarDrag();
      }}
      class={`group absolute cursor-grab active:cursor-grabbing transition-shadow ${menuOpen() ? "z-30 overflow-visible" : "z-10 overflow-hidden"
        } ${props.isContinuation && props.continuesNextDay
          ? "rounded-none border-t-0 border-b-0"
          : props.isContinuation
            ? "rounded-b-md rounded-t-none border-t-0"
            : props.continuesNextDay
              ? "rounded-t-md rounded-b-none border-b-0"
              : "rounded-md"
        } hover:z-25 hover:shadow-md select-none flex flex-col justify-between`}
      style={{
        top: `${top()}px`,
        height: `${height()}px`,
        left: `calc(${(overlapIndex() * 100) / overlapCount()}% + 2px)`,
        width: `calc(${100 / overlapCount()}% - 4px)`,
        background: props.color.bg,
        "border-left": `3px solid ${props.color.border}`,
        color: props.color.text,
        padding: "2px 6px",
      }}
      onClick={(e) => {
        if (!shouldIgnoreClick()) {
          e.stopPropagation();
          props.onClick(props.occurrence);
        }
      }}
      onContextMenu={(e) => {
        e.preventDefault();
        e.stopPropagation();
        setMenuOpen(true);
      }}
    >
      {/* Context Menu on right click */}
      <Show when={menuOpen()}>
        <div
          ref={(el) => (menuRef = el)}
          class="absolute right-1 top-1 z-50 min-w-28 rounded-lg border border-border bg-popover p-1 text-popover-foreground shadow-xl animate-in fade-in-50 zoom-in-95 pointer-events-auto"
          onPointerDown={(e) => e.stopPropagation()}
          onClick={(e) => e.stopPropagation()}
        >
          <button
            type="button"
            class="flex w-full items-center gap-2 rounded px-2 py-1.5 text-left text-xs hover:bg-accent cursor-pointer transition-colors"
            onClick={(e) => {
              e.stopPropagation();
              setMenuOpen(false);
              props.onClick(props.occurrence);
            }}
          >
            <Pencil class="size-3" />
            Alterar
          </button>
          <Show when={props.onOpenAttendance}>
            <button
              type="button"
              class="flex w-full items-center gap-2 rounded px-2 py-1.5 text-left text-xs hover:bg-accent cursor-pointer transition-colors"
              onClick={(e) => {
                e.stopPropagation();
                setMenuOpen(false);
                props.onOpenAttendance?.(props.occurrence);
              }}
            >
              <UserCheck class="size-3" />
              Presença
            </button>
          </Show>
          <Show when={props.program?.kind === "activity" && props.onOpenDraw}>
            <button
              type="button"
              class="flex w-full items-center gap-2 rounded px-2 py-1.5 text-left text-xs hover:bg-accent cursor-pointer transition-colors"
              onClick={(e) => {
                e.stopPropagation();
                setMenuOpen(false);
                props.onOpenDraw?.(props.occurrence);
              }}
            >
              <Gift class="size-3" />
              Sortear
            </button>
          </Show>
          <Show when={props.onDelete}>
            <button
              type="button"
              class="flex w-full items-center gap-2 rounded px-2 py-1.5 text-left text-xs text-destructive hover:bg-destructive/10 cursor-pointer transition-colors"
              onClick={(e) => {
                e.stopPropagation();
                setMenuOpen(false);
                props.onDelete?.(props.occurrence.id);
              }}
            >
              <Trash2 class="size-3" />
              Excluir
            </button>
          </Show>
        </div>
      </Show>

      {/* Continuation indicator from previous day */}
      <Show when={props.isContinuation}>
        <div class="text-[9px] font-mono opacity-80 leading-none py-0.5 truncate text-primary font-bold">
          ↑ Continuação
        </div>
      </Show>

      {/* Card contents */}
      <div class="min-w-0 flex-1">
        <div class="text-[11px] font-semibold whitespace-nowrap overflow-hidden text-ellipsis leading-tight">
          {props.program?.name || "Atividade"}
        </div>
        <div class="text-[10px] opacity-80 whitespace-nowrap overflow-hidden text-ellipsis leading-tight">
          {formatTimeRange(props.occurrence.starts_at, props.occurrence.ends_at)}
        </div>
        <Show when={props.occurrence.max_capacity != null}>
          <div class="text-[10px] opacity-70 mt-px truncate">
            Cap: {props.occurrence.max_capacity}
          </div>
        </Show>
      </div>

      {/* Continuation indicator to next day */}
      <Show when={props.continuesNextDay}>
        <div class="text-[9px] font-mono opacity-80 leading-none py-0.5 truncate text-primary font-bold">
          Continua no dia seguinte ↓
        </div>
      </Show>
    </div>
  );
}
