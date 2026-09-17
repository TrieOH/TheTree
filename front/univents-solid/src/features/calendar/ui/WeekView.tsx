import type { JSX } from "@solidjs/web";
import { For, Show, createMemo, createSignal, onCleanup } from "solid-js";

import ClockIcon from "~icons/lucide/clock";

import {
  addDays,
  getNowPosition,
  isToday,
  shortDayName,
  toISODate,
} from "../lib/date";
import {
  endCalendarDrag,
  getCalendarDrag,
  shouldIgnoreClick,
} from "../lib/drag-state";
import type { CalendarDragData, EventColor, OccurrenceI, ProgramI } from "../model";
import { FALLBACK_EVENT_COLOR } from "../model";
import { DraggableEvent } from "./DraggableEvent";

const Clock = ClockIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface WeekViewProps {
  currentDate: Date;
  occurrences: OccurrenceI[];
  programs: ProgramI[];
  programColors: Record<string, EventColor>;
  onSlotClick: (dateStr: string, hour: number) => void;
  onDayClick?: (date: Date) => void;
  onOccurrenceClick: (occurrence: OccurrenceI) => void;
  onDeleteOccurrence?: (occurrenceId: string) => void;
  onDropSlot?: (dateStr: string, hour: number, minute: number, data: CalendarDragData) => void;
}

const HOURS = Array.from({ length: 24 }, (_, i) => i);

function DayColumn(props: {
  day: Date;
  occurrences: OccurrenceI[];
  programs: ProgramI[];
  programColors: Record<string, EventColor>;
  nowMins: () => number;
  onSlotClick: (dateStr: string, hour: number) => void;
  onOccurrenceClick: (occurrence: OccurrenceI) => void;
  onDeleteOccurrence?: (occurrenceId: string) => void;
  onDropSlot?: (dateStr: string, hour: number, minute: number, data: CalendarDragData) => void;
}): JSX.Element {
  const [dragOverMins, setDragOverMins] = createSignal<number | null>(null);
  const dayIsToday = () => isToday(props.day);
  const dateStr = () => toISODate(props.day);

  // Memoize day occurrences so <For> only updates when data actually changes
  const dayOccurrences = createMemo(() => {
    const dayStart = new Date(
      props.day.getFullYear(),
      props.day.getMonth(),
      props.day.getDate(),
      0, 0, 0, 0,
    );
    const dayEnd = new Date(dayStart);
    dayEnd.setDate(dayEnd.getDate() + 1);

    const items = props.occurrences
      .filter((oc) => {
        const start = new Date(oc.starts_at);
        const end = new Date(oc.ends_at);
        return start < dayEnd && end > dayStart;
      })
      .map((oc) => {
        const origStart = new Date(oc.starts_at);
        const origEnd = new Date(oc.ends_at);
        const isContinuation = origStart < dayStart;
        const continuesNextDay = origEnd > dayEnd;

        return {
          occurrence: {
            ...oc,
            starts_at: isContinuation ? dayStart.toISOString() : oc.starts_at,
            ends_at: continuesNextDay ? dayEnd.toISOString() : oc.ends_at,
          },
          isContinuation,
          continuesNextDay,
        };
      });

    // Sort by starts_at
    items.sort(
      (a, b) =>
        new Date(a.occurrence.starts_at).getTime() -
        new Date(b.occurrence.starts_at).getTime(),
    );

    // Compute overlap clustering
    return items.map((item) => {
      const itemStart = new Date(item.occurrence.starts_at).getTime();
      const itemEnd = new Date(item.occurrence.ends_at).getTime();

      const overlapping = items.filter((other) => {
        const otherStart = new Date(other.occurrence.starts_at).getTime();
        const otherEnd = new Date(other.occurrence.ends_at).getTime();
        return otherStart < itemEnd && otherEnd > itemStart;
      });

      const overlapCount = Math.max(1, overlapping.length);
      const overlapIndex = overlapping.indexOf(item);

      return {
        ...item,
        overlapIndex: overlapIndex >= 0 ? overlapIndex : 0,
        overlapCount,
      };
    });
  });

  const handleDragOver = (e: DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.dataTransfer) {
      e.dataTransfer.dropEffect = "copy";
    }
    const rect = (e.currentTarget as HTMLElement).getBoundingClientRect();
    const offsetY = e.clientY - rect.top;
    const clampedY = Math.max(0, Math.min(24 * 60 - 15, offsetY));
    const snapMins = Math.floor(clampedY / 15) * 15;
    setDragOverMins(snapMins);
  };

  const handleDragLeave = (e: DragEvent) => {
    const current = e.currentTarget as HTMLElement;
    const related = e.relatedTarget as Node | null;
    if (!related || !current.contains(related)) {
      setDragOverMins(null);
    }
  };

  const handleDrop = (e: DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setDragOverMins(null);
    const data = getCalendarDrag(e);
    if (data) {
      const rect = (e.currentTarget as HTMLElement).getBoundingClientRect();
      const offsetY = e.clientY - rect.top;
      const clampedY = Math.max(0, Math.min(24 * 60 - 15, offsetY));
      const totalMinutes = Math.floor(clampedY / 15) * 15;
      const hour = Math.floor(totalMinutes / 60);
      const minute = totalMinutes % 60;
      props.onDropSlot?.(dateStr(), hour, minute, data);
    }
    endCalendarDrag();
  };

  return (
    <div
      class="min-w-30 flex-1 border-r border-border relative h-360 select-none"
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
    >
      {/* 24 Hour slot backgrounds */}
      <For each={HOURS}>
        {(h) => (
          <div
            class="h-15 border-b border-border/40 hover:bg-muted/15 transition-colors cursor-pointer box-border"
            onClick={() => {
              if (!shouldIgnoreClick()) {
                props.onSlotClick(dateStr(), h);
              }
            }}
          />
        )}
      </For>

      {/* Ghost Drop Preview */}
      <Show when={dragOverMins() !== null}>
        <div
          class="absolute left-1 right-1 rounded-md bg-primary/20 border-2 border-primary/50 border-dashed pointer-events-none z-15 transition-all flex items-center justify-center text-[10px] font-mono font-semibold text-primary"
          style={{
            top: `${dragOverMins()}px`,
            height: "60px",
          }}
        >
          {String(Math.floor((dragOverMins() ?? 0) / 60)).padStart(2, "0")}:
          {String((dragOverMins() ?? 0) % 60).padStart(2, "0")}
        </div>
      </Show>

      {/* Occurrences */}
      <For each={dayOccurrences()}>
        {(item) => (
          <DraggableEvent
            occurrence={item.occurrence}
            program={props.programs.find(
              (p) => p.id === item.occurrence.program_id,
            )}
            color={
              props.programColors[item.occurrence.program_id] ??
              FALLBACK_EVENT_COLOR
            }
            onClick={props.onOccurrenceClick}
            onDelete={props.onDeleteOccurrence}
            overlapIndex={item.overlapIndex}
            overlapCount={item.overlapCount}
            isContinuation={item.isContinuation}
            continuesNextDay={item.continuesNextDay}
          />
        )}
      </For>

      {/* Red Current Time Line */}
      <Show when={dayIsToday()}>
        <div
          class="absolute left-0 right-0 h-0.5 z-20 pointer-events-none"
          style={{
            top: `${props.nowMins()}px`,
            background: "var(--destructive, #ef4444)",
          }}
        >
          <div
            class="absolute -left-1 -top-1 size-2.5 rounded-full"
            style={{ background: "var(--destructive, #ef4444)" }}
          />
        </div>
      </Show>
    </div>
  );
}

export function WeekView(props: WeekViewProps): JSX.Element {
  const [nowMins, setNowMins] = createSignal(getNowPosition());

  const timer = setInterval(() => {
    setNowMins(getNowPosition());
  }, 30000);
  onCleanup(() => clearInterval(timer));

  const days = createMemo(() => {
    const d = new Date(props.currentDate);
    const dayOfWeek = d.getDay();
    const sunday = new Date(d);
    sunday.setDate(d.getDate() - dayOfWeek);
    sunday.setHours(0, 0, 0, 0);

    return Array.from({ length: 7 }, (_, i) => addDays(sunday, i));
  });

  const timezoneLabel = () => {
    const offset = -new Date().getTimezoneOffset() / 60;
    const sign = offset >= 0 ? "+" : "";
    return `GMT${sign}${offset}`;
  };

  const formatHourLabel = (h: number): string => `${String(h).padStart(2, "0")}:00`;

  const formatCurrentTime = (): string => {
    const now = new Date();
    const h = String(now.getHours()).padStart(2, "0");
    const m = String(now.getMinutes()).padStart(2, "0");
    return `${h}:${m}`;
  };

  return (
    <div class="calendar-scroll min-w-0 flex-1 overflow-auto bg-background">
      <div class="min-w-full">
        {/* UNIFIED STICKY HEADER ROW (SOLID OPAQUE BG, CONTINUOUS BOTTOM BORDER) */}
        <div class="sticky top-0 z-30 flex min-w-full border-b border-border bg-card select-none shadow-xs">
          {/* Timezone corner cell */}
          <div
            class="w-16 shrink-0 h-12 border-r border-border flex items-center justify-center gap-1 px-1 bg-card text-center"
            title={Intl.DateTimeFormat().resolvedOptions().timeZone}
          >
            <Clock class="size-3 text-muted-foreground/70" />
            <span class="text-[10px] font-mono text-muted-foreground font-medium">
              {timezoneLabel()}
            </span>
          </div>

          {/* 7 Days Headers */}
          <div class="grid grid-cols-7 flex-1">
            <For each={days()}>
              {(day) => {
                const dayIsToday = () => isToday(day);

                return (
                  <div
                    class={`min-w-30 flex-1 border-r border-border flex h-12 items-center justify-center gap-2 cursor-pointer select-none transition-colors hover:bg-muted/40 ${dayIsToday() ? "bg-primary/10" : "bg-card"
                      }`}
                    onClick={() => props.onDayClick?.(day)}
                  >
                    <span
                      class={`text-[11px] font-medium uppercase tracking-wider ${dayIsToday() ? "text-primary font-bold" : "text-muted-foreground"
                        }`}
                    >
                      {shortDayName(day)}
                    </span>
                    <span
                      class={`size-6 flex items-center justify-center rounded-full text-xs font-bold ${dayIsToday()
                          ? "bg-primary text-primary-foreground shadow-xs"
                          : "text-foreground"
                        }`}
                    >
                      {day.getDate()}
                    </span>
                  </div>
                );
              }}
            </For>
          </div>
        </div>

        {/* CALENDAR BODY (HOURLY GRID) */}
        <div class="flex min-w-full">
          {/* Left time column */}
          <div class="w-16 shrink-0 border-r border-border select-none bg-background">
            <div class="relative h-360">
              <For each={HOURS}>
                {(h) => (
                  <div class="h-15 border-b border-border/40 pr-2 pt-1.5 flex items-start justify-end box-border">
                    <span class="text-[11px] font-mono text-muted-foreground/80 leading-none select-none">
                      {formatHourLabel(h)}
                    </span>
                  </div>
                )}
              </For>

              {/* Current time pill on time axis */}
              <div
                class="absolute right-0 -translate-y-1/2 z-20 pointer-events-none"
                style={{ top: `${nowMins()}px` }}
              >
                <span class="bg-destructive text-destructive-foreground text-[10px] font-mono font-bold px-1.5 py-0.5 rounded-l shadow-xs">
                  {formatCurrentTime()}
                </span>
              </div>
            </div>
          </div>

          {/* 7 Days Columns */}
          <div class="grid grid-cols-7 flex-1">
            <For each={days()}>
              {(day) => (
                <DayColumn
                  day={day}
                  occurrences={props.occurrences}
                  programs={props.programs}
                  programColors={props.programColors}
                  nowMins={nowMins}
                  onSlotClick={props.onSlotClick}
                  onOccurrenceClick={props.onOccurrenceClick}
                  onDeleteOccurrence={props.onDeleteOccurrence}
                  onDropSlot={props.onDropSlot}
                />
              )}
            </For>
          </div>
        </div>
      </div>
    </div>
  );
}
