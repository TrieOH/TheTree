import type { JSX } from "@solidjs/web";
import { For, Show, createMemo, createSignal } from "solid-js";

import { formatTime, getMonthGrid, isSameDay, isToday, toISODate } from "../lib/date";
import {
  endCalendarDrag,
  getCalendarDrag,
  shouldIgnoreClick,
  startCalendarDrag,
} from "../lib/drag-state";
import type { CalendarDragData, EventColor, OccurrenceI, ProgramI } from "../model";
import { FALLBACK_EVENT_COLOR } from "../model";

export interface MonthViewProps {
  currentDate: Date;
  occurrences: OccurrenceI[];
  programs: ProgramI[];
  programColors: Record<string, EventColor>;
  onDateClick: (date: Date) => void;
  onOccurrenceClick: (occurrence: OccurrenceI) => void;
  onDropDay?: (date: Date, data: CalendarDragData) => void;
}

function MonthDayCell(props: {
  day: Date;
  currentDate: Date;
  occurrences: OccurrenceI[];
  programs: ProgramI[];
  programColors: Record<string, EventColor>;
  onDateClick: (date: Date) => void;
  onOccurrenceClick: (occurrence: OccurrenceI) => void;
  onDropDay?: (date: Date, data: CalendarDragData) => void;
}): JSX.Element {
  const [isOver, setIsOver] = createSignal(false);
  const dayIsToday = () => isToday(props.day);
  const isOtherMonth = () => props.day.getMonth() !== props.currentDate.getMonth();
  const visibleOccs = createMemo(() => props.occurrences.slice(0, 3));

  return (
    <div
      class={`min-h-28 border-r border-b border-border/70 p-1.5 transition-colors cursor-pointer ${
        isOver()
          ? "bg-primary/20 ring-2 ring-primary/40 ring-inset"
          : isOtherMonth()
          ? "bg-muted/10 text-muted-foreground/50 hover:bg-muted/20"
          : "bg-card text-foreground hover:bg-muted/30"
      }`}
      onClick={() => {
        if (!shouldIgnoreClick()) {
          props.onDateClick(props.day);
        }
      }}
      onDragOver={(e) => {
        e.preventDefault();
        e.stopPropagation();
        if (e.dataTransfer) {
          e.dataTransfer.dropEffect = "copy";
        }
        setIsOver(true);
      }}
      onDragEnter={(e) => {
        e.preventDefault();
        e.stopPropagation();
        setIsOver(true);
      }}
      onDragLeave={(e) => {
        e.preventDefault();
        e.stopPropagation();
        const current = e.currentTarget as HTMLElement;
        const related = e.relatedTarget as Node | null;
        if (!related || !current.contains(related)) {
          setIsOver(false);
        }
      }}
      onDrop={(e) => {
        e.preventDefault();
        e.stopPropagation();
        setIsOver(false);
        const data = getCalendarDrag(e);
        if (data) {
          props.onDropDay?.(props.day, data);
        }
        endCalendarDrag();
      }}
    >
      <div class="flex items-center justify-between select-none">
        <span
          class={`size-6 flex items-center justify-center rounded-full text-xs font-medium tabular-nums ${
            dayIsToday()
              ? "bg-primary text-primary-foreground font-bold"
              : "text-foreground"
          }`}
        >
          {props.day.getDate()}
        </span>
        <Show when={props.occurrences.length > 3}>
          <span class="text-[10px] text-muted-foreground font-medium pr-1">
            +{props.occurrences.length - 3}
          </span>
        </Show>
      </div>

      <div class="mt-1.5 flex flex-col gap-1">
        <For each={visibleOccs()}>
          {(oc) => {
            const prog = () => props.programs.find((p) => p.id === oc.program_id);
            const color = () => props.programColors[oc.program_id] ?? FALLBACK_EVENT_COLOR;

            const isStartDay = () => isSameDay(new Date(oc.starts_at), props.day);
            const isEndDay = () => isSameDay(new Date(oc.ends_at), props.day);
            const isMultiDay = () => !isSameDay(new Date(oc.starts_at), new Date(oc.ends_at));

            return (
              <div
                draggable="true"
                onDragStart={(e) => {
                  e.stopPropagation();
                  startCalendarDrag(e, {
                    type: "occurrence",
                    occurrenceId: oc.id,
                  });
                }}
                onDragEnd={() => {
                  endCalendarDrag();
                }}
                class={`w-full text-left text-[11px] px-1.5 py-0.5 truncate font-medium transition-all hover:opacity-85 shadow-2xs cursor-grab active:cursor-grabbing flex items-center justify-between ${
                  !isMultiDay()
                    ? "rounded"
                    : isStartDay()
                    ? "rounded-l-md rounded-r-none -mr-1.5 pr-2"
                    : isEndDay()
                    ? "rounded-r-md rounded-l-none -ml-1.5 pl-2"
                    : "rounded-none -mx-1.5 px-2"
                }`}
                style={{
                  background: color().bg,
                  color: color().text,
                  "border-left": !isMultiDay() || isStartDay() ? `3px solid ${color().border}` : "none",
                  "border-right": !isMultiDay() || isEndDay() ? `1px solid ${color().border}` : "none",
                }}
                onClick={(e) => {
                  e.stopPropagation();
                  if (!shouldIgnoreClick()) {
                    props.onOccurrenceClick(oc);
                  }
                }}
              >
                <div class="truncate flex items-center gap-1">
                  <Show when={isMultiDay() && !isStartDay()}>
                    <span class="text-[9px] font-mono text-primary font-bold">←</span>
                  </Show>
                  <Show when={isStartDay()}>
                    <span class="font-mono text-[10px] opacity-75 shrink-0">
                      {formatTime(oc.starts_at)}
                    </span>
                  </Show>
                  <span class="truncate">{prog()?.name || "Atividade"}</span>
                </div>
                <Show when={isMultiDay() && !isEndDay()}>
                  <span class="text-[9px] font-mono text-primary font-bold pl-1">→</span>
                </Show>
              </div>
            );
          }}
        </For>
      </div>
    </div>
  );
}

export function MonthView(props: MonthViewProps): JSX.Element {
  const weeks = createMemo(() => getMonthGrid(props.currentDate));
  const dowLabels = ["Dom", "Seg", "Ter", "Qua", "Qui", "Sex", "Sáb"];

  // Pre-calculate occurrences by day for O(1) cell queries
  const occurrencesByDate = createMemo(() => {
    const map = new Map<string, OccurrenceI[]>();
    for (const oc of props.occurrences) {
      const key = toISODate(new Date(oc.starts_at));
      const list = map.get(key);
      if (list) {
        list.push(oc);
      } else {
        map.set(key, [oc]);
      }
    }
    return map;
  });

  return (
    <div class="flex-1 overflow-auto p-4 bg-background">
      <div class="grid grid-cols-7 border-t border-l border-border/70 rounded-lg overflow-hidden select-none">
        <For each={dowLabels}>
          {(label) => (
            <div class="text-center py-2 text-xs font-semibold text-muted-foreground border-r border-b border-border/70 bg-muted/20">
              {label}
            </div>
          )}
        </For>
        <For each={weeks().flat()}>
          {(day) => {
            const dayKey = toISODate(day);
            const dayOccurrences = () => occurrencesByDate().get(dayKey) ?? [];

            return (
              <MonthDayCell
                day={day}
                currentDate={props.currentDate}
                occurrences={dayOccurrences()}
                programs={props.programs}
                programColors={props.programColors}
                onDateClick={props.onDateClick}
                onOccurrenceClick={props.onOccurrenceClick}
                onDropDay={props.onDropDay}
              />
            );
          }}
        </For>
      </div>
    </div>
  );
}
