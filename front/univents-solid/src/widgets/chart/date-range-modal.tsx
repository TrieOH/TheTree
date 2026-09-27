import type { JSX } from "@solidjs/web";
import {
  For,
  Show,
  createEffect,
  createMemo,
  createSignal,
  untrack,
} from "solid-js";
import ChevronLeftIcon from "~icons/lucide/chevron-left";
import ChevronRightIcon from "~icons/lucide/chevron-right";
import XIcon from "~icons/lucide/x";
import { Combobox, type ComboboxOption } from "@/shared/ui/Combobox";
import {
  addMonths,
  buildCalendarWeeks,
  isSameDay,
  isWithinRange,
} from "./date-utils";
import type { CustomRange, RangeKey, RangeOption } from "./types";

type IconComp = (props: { class?: string }) => JSX.Element;
const LucideChevronLeft = ChevronLeftIcon as unknown as IconComp;
const LucideChevronRight = ChevronRightIcon as unknown as IconComp;
const LucideClose = XIcon as unknown as IconComp;

export const PRESETS: RangeOption[] = [
  { key: "7d", label: "Últimos 7 dias", days: 7 },
  { key: "30d", label: "Últimos 30 dias", days: 30 },
  { key: "90d", label: "Últimos 90 dias", days: 90 },
  { key: "all", label: "Tudo", days: null },
];

const WEEKDAY_LABELS = ["D", "S", "T", "Q", "Q", "S", "S"];

const MONTH_LABELS = Array.from({ length: 12 }, (_, i) =>
  new Intl.DateTimeFormat("pt-BR", { month: "long" }).format(
    new Date(2000, i, 1),
  ),
);

const MONTH_OPTIONS: ComboboxOption[] = MONTH_LABELS.map((label, value) => ({
  value: String(value),
  label: label.charAt(0).toLocaleUpperCase("pt-BR") + label.slice(1),
}));

export function computeRangeWindow(
  range: RangeKey,
  customRange: CustomRange,
  now = new Date(),
): { windowStart: Date | null; windowEnd: Date | null } {
  if (range === "custom") {
    return { windowStart: customRange.from, windowEnd: customRange.to };
  }
  if (range === "all") {
    return { windowStart: null, windowEnd: null };
  }
  const preset = PRESETS.find((p) => p.key === range);
  if (!preset?.days) return { windowStart: null, windowEnd: null };
  return {
    windowStart: new Date(now.getTime() - preset.days * 86_400_000),
    windowEnd: now,
  };
}

const defaultLabelFormatter = (date: Date) =>
  new Intl.DateTimeFormat("pt-BR", {
    day: "2-digit",
    month: "short",
    year: "numeric",
  }).format(date);

interface DateRangeModalProps {
  open: boolean;
  onClose: () => void;
  range: RangeKey;
  customRange: CustomRange;
  onApply: (range: RangeKey, customRange: CustomRange) => void;
  dateFormatter?: (date: Date) => string;
}

export function DateRangeModal(props: DateRangeModalProps): JSX.Element {
  const formatLabel = () => props.dateFormatter ?? defaultLabelFormatter;

  const [draftPreset, setDraftPreset] = createSignal<RangeKey>(
    untrack(() => props.range),
  );
  const [draftFrom, setDraftFrom] = createSignal<Date | null>(
    untrack(() => props.customRange.from),
  );
  const [draftTo, setDraftTo] = createSignal<Date | null>(
    untrack(() => props.customRange.to),
  );
  const [visibleMonth, setVisibleMonth] = createSignal<Date>(
    untrack(() => props.customRange.to ?? props.customRange.from ?? new Date()),
  );
  const [hoverDate, setHoverDate] = createSignal<Date | null>(null);

  createEffect(
    () =>
      props.open
        ? { range: props.range, customRange: props.customRange }
        : null,
    (data) => {
      if (data) {
        const { windowStart, windowEnd } = computeRangeWindow(
          data.range,
          data.customRange,
        );
        setDraftPreset(data.range);
        setDraftFrom(windowStart);
        setDraftTo(windowEnd);
        setVisibleMonth(windowEnd ?? windowStart ?? new Date());
      }
    },
  );

  createEffect(
    () => props.open,
    (isOpen) => {
      if (!isOpen) return;
      const onKeyDown = (event: KeyboardEvent) => {
        if (event.key === "Escape") props.onClose();
      };
      window.addEventListener("keydown", onKeyDown);
      return () => {
        window.removeEventListener("keydown", onKeyDown);
      };
    },
  );

  const yearOptions = createMemo((): ComboboxOption[] => {
    const y = visibleMonth().getFullYear();
    const list: ComboboxOption[] = [];
    for (let i = 0; i < 9; i++) {
      const year = String(y - 7 + i);
      list.push({ value: year, label: year });
    }
    return list;
  });

  const weeks = createMemo(() => buildCalendarWeeks(visibleMonth()));

  const rangeFrom = () => draftFrom();
  const rangeTo = () =>
    draftTo() ??
    (draftFrom() && hoverDate() ? hoverDate() : null);

  const handleDayClick = (day: Date) => {
    setDraftPreset("custom");
    const currentFrom = draftFrom();
    const currentTo = draftTo();
    if (!currentFrom || currentTo) {
      setDraftFrom(day);
      setDraftTo(null);
      return;
    }
    if (day.getTime() < currentFrom.getTime()) {
      setDraftTo(currentFrom);
      setDraftFrom(day);
    } else {
      setDraftTo(day);
    }
  };

  const handlePresetClick = (key: RangeKey) => {
    const { windowStart, windowEnd } = computeRangeWindow(key, {
      from: null,
      to: null,
    });
    setDraftPreset(key);
    setDraftFrom(windowStart);
    setDraftTo(windowEnd);
    setVisibleMonth(windowEnd ?? windowStart ?? new Date());
  };

  const handleApply = () => {
    if (draftPreset() === "custom") {
      const from = draftFrom();
      if (!from) return;
      props.onApply("custom", { from, to: draftTo() ?? from });
    } else {
      props.onApply(draftPreset(), { from: null, to: null });
    }
    props.onClose();
  };

  const summaryLabel = createMemo(() => {
    const f = draftFrom();
    const t = draftTo();
    const fmt = formatLabel();

    if (draftPreset() === "custom") {
      if (!f) return "Selecione a data inicial";
      return `${fmt(f)}${t ? ` – ${fmt(t)}` : " – selecione o fim"}`;
    }
    if (f && t) {
      return `${fmt(f)} – ${fmt(t)}`;
    }
    return PRESETS.find((item) => item.key === draftPreset())?.label ?? "";
  });

  const today = new Date();

  return (
    <Show when={props.open}>
      <div
        class="fixed inset-0 z-60 flex items-center justify-center bg-foreground/40 p-2 backdrop-blur-sm transition-opacity duration-150 sm:p-4"
        onClick={() => props.onClose()}
        role="presentation"
      >
        <div
          onClick={(e) => e.stopPropagation()}
          role="dialog"
          aria-modal="true"
          aria-label="Selecionar período"
          style={{ "max-height": "calc(100vh - 2rem)" }}
          class="flex min-w-75 max-w-sm flex-1 flex-col overflow-hidden rounded-lg border border-border bg-card text-card-foreground shadow-2xl transition-all duration-150 sm:max-w-2xl"
        >
          {/* Header */}
          <div class="flex shrink-0 items-center justify-between border-b border-border px-4 py-3">
            <div>
              <p class="text-sm font-semibold">Período personalizado</p>
              <p class="hidden text-xs text-muted-foreground sm:block">
                Escolha o intervalo que deseja visualizar.
              </p>
            </div>
            <button
              type="button"
              onClick={() => props.onClose()}
              aria-label="Fechar"
              class="rounded-lg p-1 text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
            >
              <LucideClose class="size-4" />
            </button>
          </div>

          <div class="flex min-h-0 flex-1 flex-col overflow-y-auto">
            {/* Presets */}
            <div class="grid shrink-0 grid-cols-2 gap-1.5 border-b border-border p-3 sm:grid-cols-6">
              <For each={PRESETS}>
                {(p) => (
                  <button
                    type="button"
                    onClick={() => handlePresetClick(p.key)}
                    class={`shrink-0 rounded-lg px-3 py-2 text-center text-xs font-medium transition-colors ${draftPreset() === p.key
                        ? "bg-primary text-primary-foreground"
                        : "text-muted-foreground hover:bg-muted hover:text-foreground"
                      }`}
                  >
                    {p.label}
                  </button>
                )}
              </For>
              <button
                type="button"
                onClick={() => setDraftPreset("custom")}
                class={`shrink-0 rounded-lg px-3 py-2 text-center text-xs font-medium transition-colors ${draftPreset() === "custom"
                    ? "bg-primary text-primary-foreground"
                    : "text-muted-foreground hover:bg-muted hover:text-foreground"
                  }`}
              >
                Personalizado
              </button>
            </div>

            {/* Calendar */}
            <div class="flex-1 p-4">
              <div class="mb-3 flex items-center justify-between gap-2">
                <button
                  type="button"
                  onClick={() => setVisibleMonth(addMonths(visibleMonth(), -1))}
                  aria-label="Mês anterior"
                  class="rounded-lg p-1.5 text-muted-foreground hover:bg-muted hover:text-foreground"
                >
                  <LucideChevronLeft class="size-4" />
                </button>

                <div class="flex items-center gap-1.5">
                  <Combobox
                    value={String(visibleMonth().getMonth())}
                    options={MONTH_OPTIONS}
                    placeholder="Mês"
                    searchPlaceholder="Buscar mês…"
                    class="w-28 sm:w-48"
                    triggerClass="h-9 font-semibold text-xs"
                    dropdownClass="w-60 max-w-none"
                    onChange={(val) =>
                      setVisibleMonth(
                        new Date(
                          visibleMonth().getFullYear(),
                          Number(val),
                          1,
                        ),
                      )
                    }
                  />

                  <Combobox
                    value={String(visibleMonth().getFullYear())}
                    options={yearOptions()}
                    placeholder="Ano"
                    searchPlaceholder="Buscar ano…"
                    class="w-20 sm:w-32"
                    triggerClass="h-9 font-semibold text-xs"
                    dropdownClass="w-40 max-w-none"
                    onChange={(val) =>
                      setVisibleMonth(
                        new Date(
                          Number(val),
                          visibleMonth().getMonth(),
                          1,
                        ),
                      )
                    }
                  />
                </div>

                <button
                  type="button"
                  onClick={() => setVisibleMonth(addMonths(visibleMonth(), 1))}
                  aria-label="Próximo mês"
                  class="rounded-lg p-1.5 text-muted-foreground hover:bg-muted hover:text-foreground"
                >
                  <LucideChevronRight class="size-4" />
                </button>
              </div>

              <div class="grid grid-cols-7 gap-y-1">
                <For each={WEEKDAY_LABELS}>
                  {(label) => (
                    <span class="pb-1 text-center text-[10px] font-medium uppercase tracking-wide text-muted-foreground">
                      {label}
                    </span>
                  )}
                </For>

                <For each={weeks()}>
                  {(week) => (
                    <For each={week}>
                      {(day, dayIndex) => {
                        if (!day) return <div />;

                        const isStartCell = () => isSameDay(day, rangeFrom());
                        const isEndCell = () =>
                          Boolean(rangeTo() && isSameDay(day, rangeTo()));
                        const isEdgeCell = () => isStartCell() || isEndCell();
                        const inRange = () => {
                          const rf = rangeFrom();
                          const rt = rangeTo();
                          return Boolean(
                            rf && rt && isWithinRange(day, rf, rt),
                          );
                        };
                        const isToday = () => isSameDay(day, today);
                        const roundLeft = () =>
                          dayIndex() === 0 || isStartCell();
                        const roundRight = () =>
                          dayIndex() === 6 || isEndCell();

                        return (
                          <div class="relative py-0.5">
                            <Show when={inRange()}>
                              <span
                                class={`absolute inset-y-0.5 bg-primary/10 ${roundLeft() ? "left-0 rounded-l-full" : "left-0"
                                  } ${roundRight()
                                    ? "right-0 rounded-r-full"
                                    : "right-0"
                                  }`}
                              />
                            </Show>
                            <button
                              type="button"
                              onClick={() => handleDayClick(day)}
                              onMouseEnter={() => setHoverDate(day)}
                              onMouseLeave={() => setHoverDate(null)}
                              class={`relative z-10 flex size-8 mx-auto items-center justify-center rounded-full text-xs font-medium transition-colors ${isEdgeCell()
                                  ? "bg-primary text-primary-foreground font-semibold shadow-xs"
                                  : isToday()
                                    ? "border border-primary font-semibold text-primary"
                                    : "text-foreground hover:bg-muted"
                                }`}
                            >
                              {day.getDate()}
                            </button>
                          </div>
                        );
                      }}
                    </For>
                  )}
                </For>
              </div>
            </div>
          </div>

          {/* Footer */}
          <div class="flex items-center justify-between border-t border-border px-4 py-3">
            <p class="truncate text-xs text-muted-foreground">
              {summaryLabel()}
            </p>
            <div class="flex items-center gap-2">
              <button
                type="button"
                onClick={() => props.onClose()}
                class="rounded-md border border-border px-3 py-1.5 text-xs font-medium text-muted-foreground hover:bg-muted hover:text-foreground"
              >
                Cancelar
              </button>
              <button
                type="button"
                onClick={handleApply}
                class="rounded-md bg-primary px-3 py-1.5 text-xs font-medium text-primary-foreground hover:bg-primary/90"
              >
                Aplicar
              </button>
            </div>
          </div>
        </div>
      </div>
    </Show>
  );
}
