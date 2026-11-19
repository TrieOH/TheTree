import { For, createEffect, createMemo, createSignal, untrack } from "solid-js";
import ChevronLeft from "~icons/lucide/chevron-left";
import ChevronRight from "~icons/lucide/chevron-right";
import { Button } from "../button";
import { cn } from "../lib/cn";
export const dayKey = (date: Date) =>
  `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, "0")}-${String(date.getDate()).padStart(2, "0")}`;
export const validDate = (date?: Date | null): date is Date =>
  date instanceof Date && Number.isFinite(date.getTime());
export const dayTime = (date: Date) =>
  new Date(date.getFullYear(), date.getMonth(), date.getDate()).getTime();
export const formatDate = (date: Date) =>
  new Intl.DateTimeFormat("pt-BR", {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
  }).format(date);
export interface CalendarProps {
  value?: Date | null;
  from?: Date | null;
  to?: Date | null;
  min?: Date;
  max?: Date;
  isDateDisabled?: (date: Date) => boolean;
  onSelect: (date: Date) => void;
}
export function Calendar(props: CalendarProps) {
  const initial = untrack(() =>
    validDate(props.value)
      ? props.value
      : validDate(props.from)
        ? props.from
        : validDate(props.min) && dayTime(props.min) > dayTime(new Date())
          ? props.min
          : validDate(props.max) && dayTime(props.max) < dayTime(new Date())
            ? props.max
            : new Date(),
  );
  const [month, setMonth] = createSignal(
    new Date(initial.getFullYear(), initial.getMonth(), 1),
  );
  const [focused, setFocused] = createSignal(dayKey(initial));
  let grid!: HTMLDivElement;
  const unavailable = (date: Date) =>
    (validDate(props.min) && dayTime(date) < dayTime(props.min)) ||
    (validDate(props.max) && dayTime(date) > dayTime(props.max)) ||
    !!props.isDateDisabled?.(date);
  const days = createMemo(() => {
    const first = month(),
      start = new Date(
        first.getFullYear(),
        first.getMonth(),
        1 - ((first.getDay() + 6) % 7),
      );
    return Array.from(
      { length: 42 },
      (_, i) =>
        new Date(start.getFullYear(), start.getMonth(), start.getDate() + i),
    );
  });
  createEffect(
    () => ({ value: props.value, from: props.from }),
    ({ value, from }) => {
      const date = validDate(value) ? value : from;
      if (validDate(date)) {
        setMonth(new Date(date.getFullYear(), date.getMonth(), 1));
        setFocused(dayKey(date));
      }
    },
  );
  const moveMonth = (step: number) => {
    const date = new Date(month().getFullYear(), month().getMonth() + step, 1);
    setMonth(date);
    setFocused(dayKey(date));
  };
  const keyboard = (event: KeyboardEvent, date: Date) => {
    let next = new Date(date);
    const offset: Record<string, number> = {
      ArrowLeft: -1,
      ArrowRight: 1,
      ArrowUp: -7,
      ArrowDown: 7,
    };
    if (event.key in offset) next.setDate(date.getDate() + offset[event.key]);
    else if (event.key === "Home")
      next.setDate(date.getDate() - ((date.getDay() + 6) % 7));
    else if (event.key === "End")
      next.setDate(date.getDate() + 6 - ((date.getDay() + 6) % 7));
    else if (event.key === "PageUp" || event.key === "PageDown") {
      const targetMonth = date.getMonth() + (event.key === "PageUp" ? -1 : 1);
      next = new Date(
        date.getFullYear(),
        targetMonth,
        Math.min(
          date.getDate(),
          new Date(date.getFullYear(), targetMonth + 1, 0).getDate(),
        ),
      );
    } else return;
    event.preventDefault();
    setMonth(new Date(next.getFullYear(), next.getMonth(), 1));
    setFocused(dayKey(next));
    queueMicrotask(() =>
      grid
        .querySelector<HTMLButtonElement>(`[data-day="${dayKey(next)}"]`)
        ?.focus(),
    );
  };
  return (
    <div class="w-full min-w-0">
      <div class="mb-3 flex items-center justify-between gap-2">
        <Button
          size="icon"
          variant="ghost"
          aria-label="Mês anterior"
          onClick={() => moveMonth(-1)}
        >
          <ChevronLeft class="size-4" />
        </Button>
        <p aria-live="polite" class="text-sm font-medium capitalize">
          {new Intl.DateTimeFormat("pt-BR", {
            month: "long",
            year: "numeric",
          }).format(month())}
        </p>
        <Button
          size="icon"
          variant="ghost"
          aria-label="Próximo mês"
          onClick={() => moveMonth(1)}
        >
          <ChevronRight class="size-4" />
        </Button>
      </div>
      <div
        role="grid"
        aria-label="Calendário"
        ref={(element) => {
          grid = element;
        }}
        class="min-w-0"
      >
        <div role="row" class="grid grid-cols-7">
          <For each={["Seg", "Ter", "Qua", "Qui", "Sex", "Sáb", "Dom"]}>
            {(day) => (
              <span
                role="columnheader"
                class="py-1 text-center text-[10px] text-muted-foreground"
              >
                {day}
              </span>
            )}
          </For>
        </div>
        <For each={[0, 1, 2, 3, 4, 5]}>
          {(week) => (
            <div role="row" class="grid grid-cols-7">
              <For each={days().slice(week * 7, week * 7 + 7)}>
                {(date) => {
                  const selected = () =>
                    [props.value, props.from, props.to].some(
                      (d) => validDate(d) && dayKey(d) === dayKey(date),
                    );
                  const boundary = () => {
                    if (!validDate(props.from) || !validDate(props.to))
                      return undefined;
                    const time = dayTime(date),
                      from = dayTime(props.from),
                      to = dayTime(props.to);
                    if (time < from || time > to) return undefined;
                    return from === to
                      ? "single"
                      : time === from
                        ? "start"
                        : time === to
                          ? "end"
                          : "middle";
                  };
                  return (
                    <div
                      role="gridcell"
                      aria-selected={
                        selected() || boundary() ? "true" : "false"
                      }
                      data-range={boundary()}
                      class="relative isolate"
                    >
                      <span
                        aria-hidden="true"
                        class={cn(
                          "pointer-events-none absolute inset-y-0 -z-10 bg-primary/10",
                          boundary() === "start"
                            ? "left-1/2 right-0"
                            : boundary() === "end"
                              ? "left-0 right-1/2"
                              : "inset-x-0",
                          (!boundary() || boundary() === "single") && "hidden",
                          date.getDay() === 1 && "rounded-l-md",
                          date.getDay() === 0 && "rounded-r-md",
                        )}
                      />
                      <button
                        type="button"
                        data-day={dayKey(date)}
                        data-autofocus={
                          focused() === dayKey(date) ? "" : undefined
                        }
                        tabindex={focused() === dayKey(date) ? 0 : -1}
                        aria-label={new Intl.DateTimeFormat("pt-BR", {
                          dateStyle: "full",
                        }).format(date)}
                        aria-current={
                          dayKey(date) === dayKey(new Date())
                            ? "date"
                            : undefined
                        }
                        aria-disabled={unavailable(date) ? "true" : undefined}
                        class={cn(
                          "flex aspect-square w-full items-center justify-center rounded-md text-xs outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-inset",
                          date.getMonth() !== month().getMonth() &&
                            "text-muted-foreground/50",

                          selected()
                            ? "bg-primary text-primary-foreground"
                            : boundary() === "middle"
                              ? "rounded-none hover:bg-primary/15"
                              : "hover:bg-muted",
                          unavailable(date) && "cursor-not-allowed opacity-30",
                        )}
                        onFocus={() => setFocused(dayKey(date))}
                        onKeyDown={(e) => keyboard(e, date)}
                        onClick={() => {
                          if (!unavailable(date)) props.onSelect(date);
                        }}
                      >
                        {date.getDate()}
                      </button>
                    </div>
                  );
                }}
              </For>
            </div>
          )}
        </For>
      </div>
    </div>
  );
}
