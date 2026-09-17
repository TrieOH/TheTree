import type { JSX } from "@solidjs/web";
import { For, createMemo } from "solid-js";

import ChevronLeftIcon from "~icons/lucide/chevron-left";
import ChevronRightIcon from "~icons/lucide/chevron-right";

import { getMonthGrid, isToday, monthName, toISODate } from "../lib/date";

const ChevronLeft = ChevronLeftIcon as unknown as (props: { class?: string }) => JSX.Element;
const ChevronRight = ChevronRightIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface MiniCalendarProps {
  currentDate: Date;
  onDateClick: (date: Date) => void;
  onPrevMonth: () => void;
  onNextMonth: () => void;
}

export function MiniCalendar(props: MiniCalendarProps): JSX.Element {
  const weeks = createMemo(() => getMonthGrid(props.currentDate));
  const dowLabels = ["D", "S", "T", "Q", "Q", "S", "S"];
  const selectedDateStr = () => toISODate(props.currentDate);

  return (
    <div class="p-2">
      <div class="flex items-center justify-between mb-3">
        <span class="text-sm font-semibold text-foreground">
          {monthName(props.currentDate)} de {props.currentDate.getFullYear()}
        </span>
        <div class="flex gap-1">
          <button
            type="button"
            title="Mês anterior"
            class="size-7 flex items-center justify-center rounded-lg hover:bg-muted text-muted-foreground transition-colors cursor-pointer"
            onClick={() => props.onPrevMonth()}
          >
            <ChevronLeft class="size-4" />
          </button>
          <button
            type="button"
            title="Próximo mês"
            class="size-7 flex items-center justify-center rounded-lg hover:bg-muted text-muted-foreground transition-colors cursor-pointer"
            onClick={() => props.onNextMonth()}
          >
            <ChevronRight class="size-4" />
          </button>
        </div>
      </div>

      <div class="grid grid-cols-7 gap-1 text-center">
        <For each={dowLabels}>
          {(label) => (
            <div class="text-[11px] font-semibold text-muted-foreground py-1">
              {label}
            </div>
          )}
        </For>

        <For each={weeks().flat()}>
          {(day) => {
            const dayIsToday = () => isToday(day);
            const isSelected = () => toISODate(day) === selectedDateStr();
            const isOtherMonth = () => day.getMonth() !== props.currentDate.getMonth();

            return (
              <button
                type="button"
                onClick={() => props.onDateClick(day)}
                class={`size-7 mx-auto flex items-center justify-center rounded-full text-xs font-medium tabular-nums cursor-pointer transition-colors ${dayIsToday()
                  ? "bg-primary text-primary-foreground font-bold"
                  : isSelected()
                    ? "bg-primary/20 text-primary font-bold"
                    : isOtherMonth()
                      ? "text-muted-foreground/40 hover:bg-muted/50"
                      : "text-foreground hover:bg-muted"
                  }`}
              >
                {day.getDate()}
              </button>
            );
          }}
        </For>
      </div>
    </div>
  );
}
