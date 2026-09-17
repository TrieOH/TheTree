import type { JSX } from "@solidjs/web";
import { For, Show, createMemo } from "solid-js";

import { getMonthGrid, isToday, toISODate } from "../lib/date";
import type { EventColor, OccurrenceI } from "../model";

export interface YearViewProps {
  currentDate: Date;
  occurrences: OccurrenceI[];
  programColors: Record<string, EventColor>;
  onDateClick: (date: Date) => void;
}

const MONTH_NAMES = [
  "Janeiro",
  "Fevereiro",
  "Março",
  "Abril",
  "Maio",
  "Junho",
  "Julho",
  "Agosto",
  "Setembro",
  "Outubro",
  "Novembro",
  "Dezembro",
];

const DOW_LABELS = ["D", "S", "T", "Q", "Q", "S", "S"];

interface YearDayItem {
  date: Date;
  dayNum: number;
  isToday: boolean;
  isOtherMonth: boolean;
  hasEvents: boolean;
}

function MonthCard(props: {
  monthName: string;
  monthIndex: number;
  currentYear: number;
  occurrences: OccurrenceI[];
  onDateClick: (date: Date) => void;
}): JSX.Element {
  // Pre-calculate the entire month's days and event flags in a single memo
  // so individual cells don't register hundreds of separate reactive computations.
  const days = createMemo<YearDayItem[]>(() => {
    const firstDay = new Date(props.currentYear, props.monthIndex, 1);
    const grid = getMonthGrid(firstDay);

    const eventDates = new Set<string>();
    for (const oc of props.occurrences) {
      eventDates.add(toISODate(new Date(oc.starts_at)));
    }

    return grid.flat().map((day) => {
      const dateStr = toISODate(day);
      return {
        date: day,
        dayNum: day.getDate(),
        isToday: isToday(day),
        isOtherMonth: day.getMonth() !== props.monthIndex,
        hasEvents: eventDates.has(dateStr),
      };
    });
  });

  return (
    <div class="border border-border/70 rounded-xl p-3 bg-card shadow-2xs select-none">
      <h4 class="text-sm font-semibold text-foreground mb-2">
        {props.monthName}
      </h4>

      <div class="grid grid-cols-7 gap-1 text-center">
        <For each={DOW_LABELS}>
          {(label) => (
            <div class="text-[10px] font-semibold text-muted-foreground py-0.5">
              {label}
            </div>
          )}
        </For>

        <For each={days()}>
          {(day) => (
            <button
              type="button"
              onClick={() => props.onDateClick(day.date)}
              class={`relative size-7 mx-auto flex items-center justify-center rounded-full text-[11px] tabular-nums cursor-pointer transition-colors ${
                day.isToday
                  ? "bg-primary text-primary-foreground font-bold"
                  : day.hasEvents
                  ? "bg-primary/15 text-primary font-bold"
                  : day.isOtherMonth
                  ? "text-muted-foreground/30 hover:bg-muted/30"
                  : "text-foreground hover:bg-muted"
              }`}
            >
              {day.dayNum}
              <Show when={day.hasEvents && !day.isToday}>
                <span class="absolute bottom-0.5 size-1 rounded-full bg-primary" />
              </Show>
            </button>
          )}
        </For>
      </div>
    </div>
  );
}

export function YearView(props: YearViewProps): JSX.Element {
  const currentYear = () => props.currentDate.getFullYear();

  return (
    <div class="flex-1 overflow-auto p-6 bg-background">
      <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 xl:grid-cols-4 gap-6">
        <For each={MONTH_NAMES}>
          {(monthName, index) => (
            <MonthCard
              monthName={monthName}
              monthIndex={index()}
              currentYear={currentYear()}
              occurrences={props.occurrences}
              onDateClick={props.onDateClick}
            />
          )}
        </For>
      </div>
    </div>
  );
}
